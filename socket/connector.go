package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
	"time"
)

type ltvConnector struct {
	*peerProfile
	*sessionMgr

	conn net.Conn

	autoReconnectSec int // 重连间隔时间, 0为不重连

	closeSignal chan bool

	working bool // 重入锁
}

func (self *ltvConnector) Start(address string) cellnet.Peer {

	if self.working {
		return self
	}

	go self.connect(address)

	return self
}

func (self *ltvConnector) connect(address string) {
	self.working = true

	for {

		// 开始连接
		cn, err := net.Dial("tcp", address)

		// 连不上
		if err != nil {

			log.Println("[socket] cononect failed", err.Error())

			// 没重连就退出
			if self.autoReconnectSec == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}

		log.Println("connected: ", address)

		// 连上了, 记录连接
		self.conn = cn

		// 创建Session
		ses := newSession(NewPacketStream(cn), self.queue, self)
		self.sessionMgr.Add(ses)

		// 内部断开回调
		ses.OnClose = func() {
			self.sessionMgr.Remove(ses)
			self.closeSignal <- true
		}

		// 抛出事件
		self.queue.Post(NewDataEvent(Event_Connected, ses, nil))

		if <-self.closeSignal {

			self.conn = nil

			// 没重连就退出
			if self.autoReconnectSec == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue

		}

	}

	self.working = false
}

func (self *ltvConnector) Stop() {

	if self.conn != nil {
		self.conn.Close()
	}

}

func NewConnector(queue *cellnet.EvQueue) cellnet.Peer {
	return &ltvConnector{
		sessionMgr:  newSessionManager(),
		peerProfile: &peerProfile{queue: queue},
		closeSignal: make(chan bool),
	}
}
