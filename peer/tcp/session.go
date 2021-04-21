package tcp

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellpeer "github.com/davyxu/cellnet/peer"
	"github.com/davyxu/x/frame"
	"github.com/davyxu/x/io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Socket会话
type Session struct {
	xframe.PropertySet
	cellpeer.SessionIdentify

	peer   *Peer
	parent interface{}

	// Socket原始连接
	conn net.Conn

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendQueue *xframe.Pipe

	closing int64

	endNotify func()
}

// 取原始连接
func (self *Session) Raw() net.Conn {
	return self.conn
}

// 发送完队列中的封包后关闭
func (self *Session) Close() {

	closing := atomic.SwapInt64(&self.closing, 1)
	if closing != 0 {
		return
	}

	conn := self.conn

	// 关闭读
	tcpConn := conn.(*net.TCPConn)
	// 关闭读
	tcpConn.CloseRead()
	// 手动读超时
	tcpConn.SetReadDeadline(time.Now())
}

// 发送封包
func (self *Session) Send(msg interface{}) {

	// 只能通过Close关闭连接
	if msg == nil {
		return
	}

	// 已经关闭，不再发送
	if self.IsManualClosed() {
		return
	}

	// 在用户线程编码, 保证字段不会在其他线程被序列化读取
	ev := cellpeer.PackEvent(msg, &self.PropertySet)
	if ev == nil {
		return
	}
	ev.Ses = self

	self.sendQueue.Add(ev)
}

func (self *Session) IsManualClosed() bool {
	return atomic.LoadInt64(&self.closing) != 0
}

// socket层直接断开
func (self *Session) Disconnect() {
	self.conn.Close()
}

func (self *Session) readMessage() (ev *cellevent.RecvMsgEvent, err error) {

	if self.peer.Recv == nil {
		panic("no transmitter")
	}

	apply := self.peer.BeginApplyReadTimeout(self.conn)

	ev, err = self.peer.Recv(self)

	if apply {
		self.peer.EndApplyTimeout(self.conn)
	}

	return
}

// 接收循环
func (self *Session) recvLoop() {

	for self.conn != nil {

		var ev *cellevent.RecvMsgEvent
		var err error

		ev, err = self.readMessage()

		if err != nil {
			self.sendQueue.Add(nil)

			// 标记为手动关闭原因
			closedMsg := &cellevent.SessionClosed{}
			if !xio.IsEOFOrNetReadError(err) {
				closedMsg.Err = err
			}

			if self.IsManualClosed() {
				closedMsg.Reason = cellevent.CloseReason_Manual
			}

			self.peer.ProcEvent(cellevent.BuildSystemEvent(self, closedMsg))
			break
		}

		self.peer.ProcEvent(ev)
	}

	// 通知完成
	self.exitSync.Done()
}

func (self *Session) sendMessage(ev *cellevent.SendMsgEvent) (err error) {

	if self.peer.Send == nil {
		panic("no transmitter")
	}

	if self.peer.Outbound != nil {
		ev = self.peer.Outbound(ev)
	}

	apply := self.peer.BeginApplyWriteTimeout(self.conn)

	err = self.peer.Send(self, ev)

	if apply {
		self.peer.EndApplyTimeout(self.conn)
	}

	return
}

// 发送循环
func (self *Session) sendLoop() {

	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, ev := range writeList {
			self.sendMessage(ev.(*cellevent.SendMsgEvent))
		}

		if exit {
			break
		}
	}

	// 完整关闭
	self.conn.Close()

	// 通知完成
	self.exitSync.Done()
}

// 启动会话的各种资源
func (self *Session) Start() {

	atomic.StoreInt64(&self.closing, 0)

	// connector复用session时，上一次发送队列未释放可能造成问题
	self.sendQueue.Reset()

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	// 将会话添加到管理器, 在线程处理前添加到管理器(分配id), 避免ID还未分配,就开始使用id的竞态问题
	self.peer.Add(self)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.peer.Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}
	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

func newSession(conn net.Conn, p *Peer, parent interface{}) *Session {
	self := &Session{
		peer:      p,
		parent:    parent,
		conn:      conn,
		sendQueue: xframe.NewPipe(),
	}

	return self
}
