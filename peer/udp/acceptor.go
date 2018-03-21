package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net"
)

const MaxUDPRecvBuffer = 2048

type udpAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle

	localAddr *net.UDPAddr

	conn *net.UDPConn
}

func (self *udpAcceptor) Start() cellnet.Peer {

	var err error
	self.localAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.Errorf("#udp.resolve failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	self.conn, err = net.ListenUDP("udp", self.localAddr)

	if err != nil {
		log.Errorf("#udp.listen failed(%s) %s", err)
		self.SetRunning(false)
		return self
	}

	log.Infof("#udp.listen(%s) %s", self.Name(), self.Address())

	go self.accept()

	return self
}

func (self *udpAcceptor) accept() {

	self.SetRunning(true)

	for {

		buff := make([]byte, MaxUDPRecvBuffer)

		n, remoteAddr, err := self.conn.ReadFromUDP(buff)
		if err != nil {
			break
		}

		ses := newUDPSession(remoteAddr, self.conn, self)

		if n > 0 {
			ses.Recv(buff[:n])
		}

	}

	self.SetRunning(false)

}

func (self *udpAcceptor) Stop() {

	if self.conn != nil {
		self.conn.Close()
	}

	// TODO 等待accept线程结束
	self.SetRunning(false)
}

func (self *udpAcceptor) TypeName() string {
	return "udp.Acceptor"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &udpAcceptor{}

		return p
	})
}
