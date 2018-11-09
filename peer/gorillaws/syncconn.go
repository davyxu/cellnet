package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strings"
	"time"
)

type wsSyncConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *wsSession
}

func (self *wsSyncConnector) Port() int {
	if self.defaultSes.conn == nil {
		return 0
	}

	return self.defaultSes.conn.LocalAddr().(*net.TCPAddr).Port
}

func (self *wsSyncConnector) Start() cellnet.Peer {

	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = 45 * time.Second

	var finalAddress string
	if !strings.HasPrefix(self.Address(), "ws://") {
		finalAddress = "ws://" + self.Address()
	}

	conn, _, err := dialer.Dial(finalAddress, nil)

	// 发生错误时退出
	if err != nil {

		log.Debugf("#ws.connect failed(%s)@%d address: %s", self.Name(), self.defaultSes.ID(), self.Address())

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnectError{}})
		return self
	}

	self.defaultSes.conn = conn

	self.defaultSes.Start()

	self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

	return self
}

func (self *wsSyncConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *wsSyncConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *wsSyncConnector) ReconnectDuration() time.Duration {
	return 0
}

func (self *wsSyncConnector) SetReconnectDuration(v time.Duration) {

}

func (self *wsSyncConnector) Stop() {

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

}

func (self *wsSyncConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *wsSyncConnector) TypeName() string {
	return "gorillaws.SyncConnector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &wsSyncConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		self.defaultSes = newSession(nil, self, nil)

		return self
	})
}
