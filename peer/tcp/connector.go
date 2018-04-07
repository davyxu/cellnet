package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net"
	"sync"
	"time"
)

type tcpConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	defaultSes cellnet.Session

	tryConnTimes int // 尝试连接次数

	endSignal sync.WaitGroup

	reconDur time.Duration
}

func (self *tcpConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *tcpConnector) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(self.Address())

	return self
}

func (self *tcpConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *tcpConnector) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

	// 等待线程结束
	self.WaitStopFinished()

}

func (self *tcpConnector) ReconnectDuration() time.Duration {

	return self.reconDur
}

func (self *tcpConnector) SetReconnectDuration(v time.Duration) {
	self.reconDur = v
}

const reportConnectFailedLimitTimes = 3

// 连接器，传入连接地址和发送封包次数
func (self *tcpConnector) connect(address string) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		// 尝试用Socket连接地址
		conn, err := net.Dial("tcp", address)

		ses := newTCPSession(conn, self, func() {
			self.endSignal.Done()
		})
		self.defaultSes = ses

		// 发生错误时退出
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#tcp.connect failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			if self.tryConnTimes == reportConnectFailedLimitTimes {
				log.Errorf("(%s) continue reconnecting, but mute log", self.NameOrAddress())
			}

			// 没重连就退出
			if self.ReconnectDuration() == 0 {

				log.Debugf("#tcp.connect failed(%s)@%d address: %s", self.Name(), ses.ID(), self.Address())

				self.PostEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionConnectError{}})
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnectDuration())

			// 继续连接
			continue
		}

		self.endSignal.Add(1)

		ses.(interface {
			Start()
		}).Start()

		self.tryConnTimes = 0

		self.PostEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionConnected{}})

		self.endSignal.Wait()

		self.defaultSes = nil

		// 没重连就退出/主动退出
		if self.IsStopping() || self.ReconnectDuration() == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnectDuration())

		// 继续连接
		continue

	}

	self.SetRunning(false)

	self.EndStopping()
}

func (self *tcpConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *tcpConnector) TypeName() string {
	return "tcp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &tcpConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		return p
	})
}
