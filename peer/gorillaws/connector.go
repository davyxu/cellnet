package gorillaws

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

type wsConnector struct {
	peer.CoreSessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	defaultSes *wsSession

	tryConnTimes int // 尝试连接次数

	sesEndSignal sync.WaitGroup

	reconDur time.Duration
}

func (self *wsConnector) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(self.Address())

	return self
}

func (self *wsConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *wsConnector) Port() int {
	if self.defaultSes.conn == nil {
		return 0
	}

	return self.defaultSes.conn.LocalAddr().(*net.TCPAddr).Port
}

func (self *wsConnector) SetSessionManager(raw interface{}) {
	self.CoreSessionManager = raw.(peer.CoreSessionManager)
}

func (self *wsConnector) Stop() {

	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	// 通知发送关闭
	self.defaultSes.Close()

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *wsConnector) ReconnectDuration() time.Duration {

	return self.reconDur
}

func (self *wsConnector) SetReconnectDuration(v time.Duration) {
	self.reconDur = v
}

const reportConnectFailedLimitTimes = 3

func (self *wsConnector) connect(address string) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		dialer := websocket.Dialer{}
		dialer.Proxy = http.ProxyFromEnvironment
		dialer.HandshakeTimeout = 45 * time.Second

		addrObj, err := util.ParseAddress(address)
		if err != nil {
			log.Errorf("invalid address: %s", address)
			break
		}

		// 处理非法路径问题
		var finalAddress string
		if addrObj.Scheme == "ws" || addrObj.Scheme == "wss" {
			finalAddress = address
		} else {
			finalAddress = "ws://" + fmt.Sprintf("%s:%d%s", addrObj.Host, addrObj.MinPort, addrObj.Path)
		}

		conn, _, err := dialer.Dial(finalAddress, nil)
		self.defaultSes.conn = conn

		if err != nil {
			if self.tryConnTimes <= reportConnectFailedLimitTimes {

				log.Errorf("#ws.connect failed(%s) %v", self.Name(), err.Error())

				if self.tryConnTimes == reportConnectFailedLimitTimes {
					log.Errorf("(%s) continue reconnecting, but mute log", self.Name())
				}
			}

			// 没重连就退出
			if self.ReconnectDuration() == 0 || self.IsStopping() {

				self.ProcEvent(&cellnet.RecvMsgEvent{
					Ses: self.defaultSes,
					Msg: &cellnet.SessionConnectError{},
				})
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnectDuration())

			// 继续连接
			continue

		}

		self.sesEndSignal.Add(1)

		self.defaultSes.Start()

		self.tryConnTimes = 0

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

		self.sesEndSignal.Wait()

		self.defaultSes.conn = nil

		// 没重连就退出/主动退出
		if self.IsStopping() || self.ReconnectDuration() == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnectDuration())
	}

	self.SetRunning(false)
	self.EndStopping()

}

func (self *wsConnector) TypeName() string {
	return "gorillaws.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &wsConnector{}

		self.defaultSes = newSession(nil, self, func() {
			self.sesEndSignal.Done()
		})

		return self
	})
}
