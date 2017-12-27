package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
)

type udpConnector struct {
	internal.PeerShare
	localAddr *net.UDPAddr
	conn      *net.UDPConn
}

func (self *udpConnector) Start() cellnet.Peer {

	var err error
	self.localAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.Errorf("#resolve udp address failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	go self.connect()

	return self
}

func (self *udpConnector) connect() {

	var err error
	self.conn, err = net.DialUDP("udp", nil, self.localAddr)
	if err != nil {

		log.Errorf("#connect failed(%s) %v", self.NameOrAddress(), err.Error())
		return
	}

	ses := newUDPSession(self.localAddr, self.conn, &self.PeerShare, nil)

	ses.Start()

	self.CallInboundProc(&cellnet.RecvMsgEvent{ses, &comm.SessionConnected{}})

	buff := make([]byte, 4096)
	for {

		n, remoteAddr, err := self.conn.ReadFromUDP(buff)
		if err != nil {

			log.Errorln("disconnected:", remoteAddr.String())
			break
		}

		err = ses.OnRecv(buff[:n])

		if err != nil {
			break
		}
	}
}

func (self *udpConnector) IsConnector() bool {
	return true
}

func (self *udpConnector) Stop() {

	if self.conn != nil {
		self.conn.Close()
	}
}

func init() {

	cellnet.RegisterPeerCreator("udp.Connector", func() cellnet.Peer {
		p := &udpConnector{}

		p.Init(p)

		return p
	})
}
