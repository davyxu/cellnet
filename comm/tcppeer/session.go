package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"net"
	"sync"
)

// Socket会话
type tcpSession struct {
	cellnet.CoreSessionShare

	// Socket原始连接
	conn net.Conn

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendChan chan interface{}

	cleanupGuard sync.Mutex

	endNotify func()
}

// 取原始连接
func (self *tcpSession) Raw() interface{} {
	return self.conn
}

func (self *tcpSession) Close() {
	self.sendChan <- nil
}

// 发送封包
func (self *tcpSession) Send(msg interface{}) {
	self.sendChan <- msg
}

// 接收循环
func (self *tcpSession) recvLoop() {

	var readEvt cellnet.ReadStreamEvent
	readEvt.Ses = self

	for self.conn != nil {

		// 发送接收消息，要求读取数据
		raw := self.PeerShare.CallInboundProc(&readEvt)

		if err, ok := raw.(error); ok && err != nil && self.conn != nil {

			self.PeerShare.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{err.Error()}})

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
		self.PeerShare.CallOutboundProc(&cellnet.SendMsgEvent{self, msg})
	}

	self.cleanup()
}

// 清理资源
func (self *tcpSession) cleanup() {

	self.cleanupGuard.Lock()

	defer self.cleanupGuard.Unlock()

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
	self.Peer().(cellnet.SessionManager).Add(self)

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(cellnet.SessionManager).Remove(self)

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

func newTCPSession(conn net.Conn, peerShare *cellnet.CoreCommunicatePeer, endNotify func()) cellnet.Session {
	self := &tcpSession{
		conn:      conn,
		endNotify: endNotify,
		sendChan:  make(chan interface{}, SendQueueLen),
	}

	self.PeerShare = peerShare

	return self
}
