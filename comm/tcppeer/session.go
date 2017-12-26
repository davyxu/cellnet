package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
)

// Socket会话
type tcpSession struct {

	// Socket原始连接
	conn net.Conn

	tag interface{}

	// 退出同步器
	exitSync sync.WaitGroup

	// 归属的通讯端
	peer *internal.PeerShare

	id int64

	// 发送队列
	sendChan chan interface{}

	endNotify func()
}

// 取原始连接
func (self *tcpSession) Raw() interface{} {
	return self.conn
}

func (self *tcpSession) Tag() interface{} {
	return self.tag
}

func (self *tcpSession) SetTag(v interface{}) {
	self.tag = v
}

func (self *tcpSession) ID() int64 {
	return self.id
}

func (self *tcpSession) SetID(id int64) {
	self.id = id
}

func (self *tcpSession) Close() {
	self.sendChan <- nil
}

// 取会话归属的通讯端
func (self *tcpSession) Peer() cellnet.Peer {
	return self.peer.Peer()
}

// 发送封包
func (self *tcpSession) Send(msg interface{}) {
	self.sendChan <- msg
}

// 接收循环
func (self *tcpSession) recvLoop() {

	for self.conn != nil {

		// 发送接收消息，要求读取数据
		raw := self.peer.CallInboundProc(&cellnet.ReadEvent{self})

		// 连接断开
		if raw != nil && self.conn != nil {

			self.peer.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})
			//self.peer.CallInboundProc(cellnet.RecvErrorEvent{self, raw.(error)})
			break
		}
	}

	self.cleanup()
}

// 发送循环
func (self *tcpSession) sendLoop() {

	// 遍历要发送的数据
	for msg := range self.sendChan {

		// nil表示需要退出会话通讯
		if msg == nil {
			break
		}

		// 要求发送数据
		err := self.peer.CallOutboundProc(&cellnet.SendMsgEvent{self, msg})

		// 发送错误时派发事件
		if err != nil {
			self.peer.CallInboundProc(&cellnet.SendMsgErrorEvent{self, err.(error), msg})
			break
		}

	}

	self.cleanup()
}

// 清理资源
func (self *tcpSession) cleanup() {

	// 关闭连接
	if self.conn != nil {
		self.conn.Close()
		self.conn = nil
	}

	// 关闭发送队列
	if self.sendChan != nil {
		close(self.sendChan)
		self.sendChan = nil
	}

	// 通知完成
	self.exitSync.Done()
}

// 启动会话的各种资源
func (self *tcpSession) Start() {

	// 将会话添加到管理器
	self.Peer().(internal.SessionManager).Add(self)

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(internal.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}

	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

// 默认10个长度的发送队列
const SendQueueLen = 100

func newTCPSession(conn net.Conn, peer *internal.PeerShare, endNotify func()) cellnet.Session {
	return &tcpSession{
		conn:      conn,
		peer:      peer,
		endNotify: endNotify,
		sendChan:  make(chan interface{}, SendQueueLen),
	}
}
