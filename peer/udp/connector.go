package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net"
)

type udpConnector struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	remoteAddr *net.UDPAddr

	defaultSes *udpSession
}

func (self *udpConnector) Start() cellnet.Peer {

	var err error
	self.remoteAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.Errorf("#resolve udp address failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	go self.connect()

	return self
}

func (self *udpConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *udpConnector) connect() {

	conn, err := net.DialUDP("udp", nil, self.remoteAddr)
	if err != nil {

		log.Errorf("#udp.connect failed(%s) %v", self.NameOrAddress(), err.Error())
		return
	}

	self.defaultSes.conn = conn

	ses := self.defaultSes

	self.PostEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionConnected{}})

	recvBuff := make([]byte, MaxUDPRecvBuffer)

	self.SetRunning(true)

	for self.IsRunning() {

		n, _, err := conn.ReadFromUDP(recvBuff)
		if err != nil {
			break
		}

		if n > 0 {
			ses.Recv(recvBuff[:n])
		}

	}
}

func (self *udpConnector) Stop() {

	self.SetRunning(false)

	if self.defaultSes.conn != nil {
		self.defaultSes.conn.Close()
	}
}

func (self *udpConnector) TypeName() string {
	return "udp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &udpConnector{}

		p.defaultSes = &udpSession{
			pInterface:     p,
			CoreProcBundle: &p.CoreProcBundle,
		}

		return p
	})
}
