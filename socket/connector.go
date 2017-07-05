package socket

import (
	"net"
	"time"

	"github.com/davyxu/cellnet"
)

type socketConnector struct {
	*peerBase
	*SessionManager

	conn net.Conn

	autoReconnectSec int // 重连间隔时间, 0为不重连

	tryConnTimes int // 尝试连接次数

	closeSignal chan bool

	defaultSes cellnet.Session
}

// 自动重连间隔=0不重连
func (self *socketConnector) SetAutoReconnectSec(sec int) {
	self.autoReconnectSec = sec
}

func (self *socketConnector) Start(address string) cellnet.Peer {

	self.waitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(address)

	return self
}

const reportConnectFailedLimitTimes = 3

func (self *socketConnector) connect(address string) {

	self.SetRunning(true)
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

				systemError(nil, cellnet.Event_ConnectFailed, errToResult(err), self.safeRecvHandler())
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

		self.SessionManager.Add(ses)

		// 内部断开回调
		ses.OnClose = func() {
			self.SessionManager.Remove(ses)
			self.closeSignal <- true
		}

		systemEvent(ses, cellnet.Event_Connected, self.safeRecvHandler())

		if <-self.closeSignal {

			self.conn = nil

			// 没重连就退出/主动退出
			if self.isStopping() || self.autoReconnectSec == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue

		}

	}

	self.SetRunning(false)

	self.endStopping()
}

func (self *socketConnector) Stop() {

	if !self.IsRunning() {
		return
	}

	if self.isStopping() {
		return
	}

	self.startStopping()

	// socket断开, 后续触发一系列事件通知
	self.CloseAllSession()

	// 等待线程结束
	self.waitStopFinished()
}

func (self *socketConnector) DefaultSession() cellnet.Session {
	return self.defaultSes
}

func NewConnector(q cellnet.EventQueue) cellnet.Peer {

	return NewConnectorBySessionManager(q, ClientSessionManager)
}

func NewConnectorBySessionManager(q cellnet.EventQueue, sm *SessionManager) cellnet.Peer {
	self := &socketConnector{
		SessionManager: sm,
		peerBase:       newPeerBase(q),
		closeSignal:    make(chan bool),
	}

	return self
}
