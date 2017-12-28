package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
)

const MaxUDPRecvBuffer = 2048

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

	buff := make([]byte, MaxUDPRecvBuffer)

	var recentAddr addressPair
	var recentSes *udpSession

	for {

		n, remoteAddr, err := self.conn.ReadFromUDP(buff)
		if err != nil {
			break
		}

		addr := makeAddrKey(remoteAddr)

		var ses *udpSession

		if recentAddr == addr {

			ses = recentSes

		} else {

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
			}

			recentAddr = addr
			recentSes = ses
		}

		// 将数据拷贝到会话缓冲区

		if n > 0 {
			ses.OnRecv(buff[:n])

			// 并发处理封包
			go ses.ProcPacket()

		} else { // n=0情况 Mono在连上后，需要发一个包, 直接处理会发生EOF
			ses.KeepAlive()
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
