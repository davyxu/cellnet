package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"net"
	"sync"
	"sync/atomic"
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
	sendQueue *util.Pipe

	cleanupGuard sync.Mutex

	endNotify func()

	manualClosed int64
}

func (self *tcpSession) Peer() cellnet.Peer {
	return self.pInterface
}

// 取原始连接
func (self *tcpSession) Raw() interface{} {
	if self.conn == nil {
		return nil
	}

	return self.conn
}

func (self *tcpSession) Close() {
	atomic.StoreInt64(&self.manualClosed, 1)
	self.sendQueue.Add(nil)
}

// 发送封包
func (self *tcpSession) Send(msg interface{}) {

	// 只能通过Close关闭连接
	if msg == nil {
		return
	}

	self.sendQueue.Add(msg)
}

// 接收循环
func (self *tcpSession) recvLoop() {

	for self.conn != nil {

		msg, err := self.ReadMessage(self)

		if err != nil {
			if !util.IsEOFOrNetReadError(err) {
				log.Errorf("session closed, sesid: %d, err: %s", self.ID(), err)
			}

			self.sendQueue.Add(nil)

			manualClosed := atomic.LoadInt64(&self.manualClosed)

			// 标记为手动关闭原因
			closedMsg := &cellnet.SessionClosed{}
			if manualClosed != 0 {
				closedMsg.Reason = cellnet.CloseReason_Manual
			}

			self.PostEvent(&cellnet.RecvMsgEvent{self, closedMsg})
			break
		}

		self.PostEvent(&cellnet.RecvMsgEvent{self, msg})
	}

	self.cleanup()
}

// 发送循环
func (self *tcpSession) sendLoop() {

	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, msg := range writeList {

			self.SendMessage(&cellnet.SendMsgEvent{self, msg})
		}

		if exit {
			break
		}
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

	// 通知完成
	self.exitSync.Done()
}

// 启动会话的各种资源
func (self *tcpSession) Start() {

	atomic.StoreInt64(&self.manualClosed, 0)

	// connector复用session时，上一次发送队列未释放可能造成问题
	self.sendQueue.Reset()

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

func newSession(conn net.Conn, p cellnet.Peer, endNotify func()) *tcpSession {
	self := &tcpSession{
		conn:       conn,
		endNotify:  endNotify,
		sendQueue:  util.NewPipe(),
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}

	return self
}
