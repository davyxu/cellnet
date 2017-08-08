package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"net"
	"time"
)

// 连接器, 可由Peer转换
type Connector interface {

	// 连接后的Session
	DefaultSession() cellnet.Session

	// 自动重连间隔, 0表示不重连, 默认不重连
	SetAutoReconnectSec(sec int)
}

type socketConnector struct {
	*socketPeer

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
	self.SetAddress(address)

	for {

		self.tryConnTimes++

		// 开始连接
		conn, err := net.Dial("tcp", address)

		// 连不上
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#connect failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			if self.tryConnTimes == reportConnectFailedLimitTimes {
				log.Errorf("(%s) continue reconnecting, but mute log", self.NameOrAddress())
			}

			// 没重连就退出
			if self.autoReconnectSec == 0 {

				extend.PostSystemEvent(nil, cellnet.Event_ConnectFailed, self.ChainListRecv(), errToResult(err))
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}

		ses := newSession(conn, self)
		self.defaultSes = ses
		ses.run()

		self.tryConnTimes = 0

		// 创建Session

		self.Add(ses)

		// 内部断开回调
		ses.OnClose = func() {
			self.Remove(ses)
			self.closeSignal <- true
		}

		extend.PostSystemEvent(ses, cellnet.Event_Connected, self.ChainListRecv(), cellnet.Result_OK)

		if <-self.closeSignal {

			self.defaultSes = nil

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

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

	// 等待线程结束
	self.waitStopFinished()
}

func (self *socketConnector) DefaultSession() cellnet.Session {
	return self.defaultSes
}

func (self *socketConnector) RPCSession() cellnet.Session {
	return self.defaultSes
}

func NewConnector(q cellnet.EventQueue) cellnet.Peer {

	return NewConnectorBySessionManager(q, cellnet.NewSessionManager())
}

func NewConnectorBySessionManager(q cellnet.EventQueue, sm cellnet.SessionManager) cellnet.Peer {
	self := &socketConnector{
		socketPeer:  newSocketPeer(q, sm),
		closeSignal: make(chan bool),
	}

	return self
}
