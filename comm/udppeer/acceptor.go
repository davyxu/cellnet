package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
)

type udpAcceptor struct {
	internal.PeerShare
	localAddr *net.UDPAddr

	conn *net.UDPConn

	sesByAddress map[*net.UDPAddr]*udpSession
}

func (self *udpAcceptor) Start() cellnet.Peer {

	var err error
	self.localAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.Errorf("#resolve udp address failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	go self.listen()

	return self
}

func (self *udpAcceptor) listen() {

	var err error
	self.conn, err = net.ListenUDP("udp", self.localAddr)

	if err != nil {
		log.Errorln("listen failed:", self.localAddr.String())
		return
	}

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	buff := make([]byte, 4096)

	for {

		n, remoteAddr, err := self.conn.ReadFromUDP(buff)
		if err != nil {

			log.Errorln("disconnected:", err)
			//self.FireEvent(cellnet.SessionClosedEvent{nil})
			break
		}

		ses := self.sesByAddress[remoteAddr]

		if ses == nil {

			ses = newUDPSession(remoteAddr, self.conn, &self.PeerShare, nil)

			ses.Start()

			self.sesByAddress[remoteAddr] = ses
		}

		ses.OnRecv(buff[:n])

	}

}

func (self *udpAcceptor) IsAcceptor() bool {
	return true
}

func (self *udpAcceptor) Stop() {

	self.conn.Close()
}

func init() {

	cellnet.RegisterPeerCreator("udp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &udpAcceptor{
			sesByAddress: make(map[*net.UDPAddr]*udpSession),
		}

		p.Init(p, config)

		return p
	})
}
