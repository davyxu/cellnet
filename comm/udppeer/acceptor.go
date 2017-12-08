package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
)

type v4Address struct {
	IP   [4]byte
	Port int
}

type v6Address struct {
	IP   [16]byte
	Port int
}

func makeAddrKey(addr *net.UDPAddr) interface{} {

	switch len(addr.IP) {
	case net.IPv4len:
		var ret v4Address
		for i := 0; i < net.IPv4len; i++ {
			ret.IP[i] = addr.IP[i]
		}
		ret.Port = addr.Port

		return ret
	case net.IPv6len:
		var ret v6Address
		for i := 0; i < net.IPv6len; i++ {
			ret.IP[i] = addr.IP[i]
		}
		ret.Port = addr.Port

		return ret
	}

	return nil
}

type udpAcceptor struct {
	internal.PeerShare
	localAddr *net.UDPAddr

	conn *net.UDPConn

	sesByAddress map[interface{}]*udpSession
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

			log.Errorln("read error:", err)
			//self.FireEvent(cellnet.SessionClosedEvent{nil})
			break
		}

		addr := makeAddrKey(remoteAddr)

		ses := self.sesByAddress[addr]

		if ses == nil {

			ses = newUDPSession(remoteAddr, self.conn, &self.PeerShare, nil)

			ses.Start()

			self.sesByAddress[addr] = ses

			self.FireEvent(cellnet.SessionAcceptedEvent{ses})
		}

		log.Debugln("recv data", buff[:n])

		err = ses.OnRecv(buff[:n])

		if err != nil {
			delete(self.sesByAddress, addr)
		}
	}

}

func (self *udpAcceptor) IsAcceptor() bool {
	return true
}

func (self *udpAcceptor) Stop() {

	if self.conn != nil {
		self.conn.Close()
	}

}

func init() {

	cellnet.RegisterPeerCreator("udp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &udpAcceptor{
			sesByAddress: make(map[interface{}]*udpSession),
		}

		p.Init(p, config)

		return p
	})
}
