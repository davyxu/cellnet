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
	internal.CommunicatePeer
	internal.PeerInfo
	localAddr *net.UDPAddr

	conn *net.UDPConn

	sesByAddress sync.Map

	pktPool sync.Pool
}

func (self *udpAcceptor) Start() cellnet.Peer {

	var err error
	self.localAddr, err = net.ResolveUDPAddr("udp", self.Address())

	if err != nil {

		log.Errorf("#resolve udp address failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	self.conn, err = net.ListenUDP("udp", self.localAddr)

	if err != nil {
		log.Errorln("listen failed:", err)
		return self
	}

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go self.accept()

	return self
}

func (self *udpAcceptor) accept() {

	var recentAddr addressPair
	var recentSes *udpSession

	for {

		buff := self.pktPool.Get().([]byte)[:MaxUDPRecvBuffer]

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

				ses = newUDPSession(remoteAddr, self.conn, &self.CommunicatePeer, func() {
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
			ses.Recv(buff[:n])

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

func (self *udpAcceptor) TypeName() string {
	return "udp.Acceptor"
}

func init() {

	cellnet.RegisterPeerCreator(func() cellnet.Peer {
		p := &udpAcceptor{}

		p.pktPool.New = func() interface{} {
			return make([]byte, MaxUDPRecvBuffer)
		}

		p.Init(p)

		return p
	})
}
