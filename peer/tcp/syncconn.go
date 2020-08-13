package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/ulog"
	"net"
	"time"
)

type tcpSyncConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *tcpSession
}

func (self *tcpSyncConnector) Port() int {
	conn := self.defaultSes.Conn()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.TCPAddr).Port
}

func (self *tcpSyncConnector) Start() cellnet.Peer {

	// 尝试用Socket连接地址
	conn, err := net.Dial("tcp", self.Address())

	// 发生错误时退出
	if err != nil {

		ulog.Debugf("#tcp.connect failed(%s)@%d address: %s", self.Name(), self.defaultSes.ID(), self.Address())

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnectError{}})
		return self
	}

	self.defaultSes.setConn(conn)

	self.ApplySocketOption(conn)

	self.defaultSes.Start()

	self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

	return self
}

func (self *tcpSyncConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *tcpSyncConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *tcpSyncConnector) ReconnectInterval() time.Duration {
	return 0
}

func (self *tcpSyncConnector) SetReconnectInterval(v time.Duration) {

}

func (self *tcpSyncConnector) Stop() {

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

}

func (self *tcpSyncConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *tcpSyncConnector) TypeName() string {
	return "tcp.SyncConnector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &tcpSyncConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		self.defaultSes = newSession(nil, self, nil)

		self.CoreTCPSocketOption.Init()

		return self
	})
}
