package socket

import (
	"net"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type socketConnector struct {
	*peerBase
	*sessionMgr

	conn net.Conn

	autoReconnectSec int // 重连间隔时间, 0为不重连

	closeSignal chan bool

	working bool // 重入锁

	defaultSes cellnet.Session

	tryConnTimes int // 尝试连接次数
}

// 自动重连间隔=0不重连
func (self *socketConnector) SetAutoReconnectSec(sec int) {
	self.autoReconnectSec = sec
}

func (self *socketConnector) Start(address string) cellnet.Peer {

	if self.working {
		return self
	}

	go self.connect(address)

	return self
}

const reportConnectFailedLimitTimes = 3

func (self *socketConnector) connect(address string) {
	self.working = true
	self.address = address

	for {

		self.tryConnTimes++

		// 开始连接
		conn, err := net.Dial("tcp", address)

		// 连不上
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#connect failed(%s) %v", self.nameOrAddress(), err.Error())
			}

			if self.tryConnTimes == reportConnectFailedLimitTimes {
				log.Errorf("(%s) continue reconnecting, but mute log", self.nameOrAddress())
			}

			// 没重连就退出
			if self.autoReconnectSec == 0 {

				callSystemEvent(cellnet.SessionEvent_ConnectFailed, &coredef.SessionConnectFailed{Reason: err.Error()}, self.safeRecvHandler())
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}

		ses := newSession(self.genPacketStream(conn), self)
		self.defaultSes = ses
		ses.run()

		self.tryConnTimes = 0

		// 连上了, 记录连接
		self.conn = conn

		// 创建Session

		self.sessionMgr.Add(ses)

		// 内部断开回调
		ses.OnClose = func() {
			self.sessionMgr.Remove(ses)
			self.closeSignal <- true
		}

		callSystemEventByMeta(ses, cellnet.SessionEvent_Connected, Meta_SessionConnected, self.safeRecvHandler())

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

func (self *socketConnector) Stop() {

	if self.conn != nil {

		self.conn.Close()
	}

}

func (self *socketConnector) DefaultSession() cellnet.Session {
	return self.defaultSes
}

func NewConnector(evq cellnet.EventQueue) cellnet.Peer {
	self := &socketConnector{
		sessionMgr:  newSessionManager(),
		peerBase:    newPeerBase(evq),
		closeSignal: make(chan bool),
	}

	return self
}
