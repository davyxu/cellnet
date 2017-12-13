package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
	"time"
)

type Connector interface {
	SetAutoReconnectSec(sec int)
}

type tcpConnector struct {
	internal.PeerShare

	defaultSes cellnet.Session

	autoReconnectSec int // 重连间隔时间, 0为不重连

	tryConnTimes int // 尝试连接次数

	endSignal sync.WaitGroup
}

func (self *tcpConnector) IsConnector() bool {
	return true
}

// 自动重连间隔=0不重连
func (self *tcpConnector) SetAutoReconnectSec(sec int) {
	self.autoReconnectSec = sec
}

func (self *tcpConnector) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(self.PeerAddress)

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

const reportConnectFailedLimitTimes = 3

// 连接器，传入连接地址和发送封包次数
func (self *tcpConnector) connect(address string) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		// 尝试用Socket连接地址
		conn, err := net.Dial("tcp", address)

		ses := newTCPSession(conn, &self.PeerShare, func() {
			self.endSignal.Done()
		})
		self.defaultSes = ses

		// 发生错误时退出
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#connect failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			if self.tryConnTimes == reportConnectFailedLimitTimes {
				log.Errorf("(%s) continue reconnecting, but mute log", self.NameOrAddress())
			}

			// 没重连就退出
			if self.autoReconnectSec == 0 {

				self.FireEvent(&cellnet.SessionConnectErrorEvent{ses, err})
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}

		self.endSignal.Add(1)

		ses.(interface {
			Start()
		}).Start()

		self.tryConnTimes = 0

		self.FireEvent(&cellnet.SessionConnectedEvent{ses})

		self.endSignal.Wait()

		self.defaultSes = nil

		// 没重连就退出/主动退出
		if self.IsStopping() || self.autoReconnectSec == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

		// 继续连接
		continue

	}

	self.SetRunning(false)

	self.EndStopping()
}

func init() {

	cellnet.RegisterPeerCreator("tcp.Connector", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &tcpConnector{}

		p.Init(p, config)

		return p
	})
}
