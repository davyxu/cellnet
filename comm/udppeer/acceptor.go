package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
)

type udpAcceptor struct {
	internal.PeerShare
	localAddr *net.UDPAddr

	conn *net.UDPConn

	sesByAddress sync.Map
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
			//self.CallInboundProc(cellnet.SessionClosedEvent{nil})
			break
		}

		addr := makeAddrKey(remoteAddr)

		var ses *udpSession

		raw, ok := self.sesByAddress.Load(addr)

		if ok {

			ses = raw.(*udpSession)

		} else {

			ses = newUDPSession(remoteAddr, self.conn, &self.PeerShare, func() {
				self.removeAddress(addr)
			})

			ses.Start()

			self.sesByAddress.Store(addr, ses)

			self.CallInboundProc(&cellnet.RecvMsgEvent{ses, &comm.SessionAccepted{}})

			// mono首次封包是空
			if n == 0 {
				ses.HeartBeat()

				continue
			}
		}

		err = ses.OnRecv(buff[:n])

		if err != nil {
			self.removeAddress(addr)
		}
	}

}

func (self *udpAcceptor) removeAddress(pair addressPair) {

	self.sesByAddress.Delete(pair)
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

	cellnet.RegisterPeerCreator("udp.Acceptor", func() cellnet.Peer {
		p := &udpAcceptor{}

		p.Init(p)

		return p
	})
}
