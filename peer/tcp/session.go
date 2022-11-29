package tcp

import (
	"fmt"
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellpeer "github.com/davyxu/cellnet/peer"
	"github.com/davyxu/x/container"
	"github.com/davyxu/x/io"
	xnet "github.com/davyxu/x/net"
	"github.com/davyxu/xlog"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Socket会话
type Session struct {
	xcontainer.Mapper
	cellpeer.SessionIdentify

	Peer   *Peer
	parent any

	// Socket原始连接
	conn net.Conn

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendQueue *xcontainer.Pipe

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
func (self *Session) Send(msg any) {

	// 只能通过Close关闭连接
	if msg == nil {
		return
	}

	// 已经关闭，不再发送
	if self.IsManualClosed() {
		return
	}

	// 在用户线程编码, 保证字段不会在其他线程被序列化读取
	ev := cellpeer.PackEvent(msg, &self.Mapper)
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

func (self *Session) readMessage() (ev *cellevent.RecvMsg, err error) {

	if self.Peer.OnRecv == nil {
		panic("peer.OnRecv not set")
	}

	apply := self.Peer.BeginApplyReadTimeout(self.conn)

	self.Peer.ProctectCall(func() {

		ev, err = self.Peer.OnRecv(self)

	}, func(raw any) {
		var ok bool
		if err, ok = raw.(error); !ok {
			err = fmt.Errorf("recv panic: %+v", raw)
		}
	})

	if apply {
		self.Peer.EndApplyTimeout(self.conn)
	}

	return
}

// 接收循环
func (self *Session) recvLoop() {

	for {

		var ev *cellevent.RecvMsg
		var err error

		ev, err = self.readMessage()

		if err != nil {
			self.sendQueue.Stop(false)

			// 标记为手动关闭原因
			closedMsg := &cellevent.SessionClosed{}
			if !xio.IsEOFOrNetReadError(err) {
				closedMsg.Err = err
			}

			if self.IsManualClosed() {
				closedMsg.Reason = cellevent.CloseReason_Manual
			}

			self.Peer.ProcEvent(cellevent.BuildSystemEvent(self, closedMsg))
			break
		}

		self.Peer.ProcEvent(ev)
	}

	// 通知完成
	self.exitSync.Done()
}

var (
	OnSendCrash = func(raw any) {
		xlog.Errorf("send panic: %+v", raw)
	}
)

func (self *Session) sendMessage(ev *cellevent.SendMsg) (err error) {

	if self.Peer.OnSend == nil {
		panic("peer.OnSend not set")
	}

	if self.Peer.OnOutbound != nil {
		ev = self.Peer.OnOutbound(ev)
	}

	apply := self.Peer.BeginApplyWriteTimeout(self.conn)

	self.Peer.ProctectCall(func() {

		err = self.Peer.OnSend(self, ev)

	}, OnSendCrash)

	if apply {
		self.Peer.EndApplyTimeout(self.conn)
	}

	return
}

// 发送循环
func (self *Session) sendLoop() {

	var writeList []any

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, ev := range writeList {
			self.sendMessage(ev.(*cellevent.SendMsg))
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
	self.Peer.AddSession(self)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer.RemoveSession(self)

		if self.endNotify != nil {
			self.endNotify()
		}
	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

func newSession(conn net.Conn, p *Peer, parent any) *Session {
	self := &Session{
		Peer:      p,
		parent:    parent,
		conn:      conn,
		sendQueue: xcontainer.NewPipe(),
	}

	return self
}

// 获取session远程的地址
func GetRemoteAddrss(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if rawSes, ok := ses.(*Session); ok {
		return rawSes.Raw().RemoteAddr().String()
	}

	return ""
}

func GetRemoteHost(ses cellnet.Session) string {
	addr := GetRemoteAddrss(ses)
	host, _, err := xnet.SpliteAddress(addr)
	if err == nil {
		return host
	}

	return ""
}
