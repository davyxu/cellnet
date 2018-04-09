package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"io"
	"net"
	"sync"
)

// Socket会话
type tcpSession struct {
	peer.CoreContextSet
	peer.CoreSessionIdentify
	*peer.CoreProcBundle

	pInterface cellnet.Peer

	// Socket原始连接
	conn net.Conn

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendChan chan interface{}

	cleanupGuard sync.Mutex

	endNotify func()
}

func (self *tcpSession) Peer() cellnet.Peer {
	return self.pInterface
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

func isEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

// 接收循环
func (self *tcpSession) recvLoop() {

	for self.conn != nil {

		msg, err := self.ReadMessage(self)

		if err != nil {
			if !isEOFOrNetReadError(err) {
				log.Errorln("session closed:", err)
			}

			self.PostEvent(&cellnet.RecvMsgEvent{self, &cellnet.SessionClosed{}})
			break
		}

		self.PostEvent(&cellnet.RecvMsgEvent{self, msg})
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

		// TODO SendMsgEvent并不是很有意义
		self.SendMessage(&cellnet.SendMsgEvent{self, msg})

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
	self.Peer().(peer.SessionManager).Add(self)

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(peer.SessionManager).Remove(self)

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

func newTCPSession(conn net.Conn, p cellnet.Peer, endNotify func()) cellnet.Session {
	self := &tcpSession{
		conn:       conn,
		endNotify:  endNotify,
		sendChan:   make(chan interface{}, SendQueueLen),
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}

	return self
}
