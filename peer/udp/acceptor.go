package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"net"
	"time"
)

const MaxUDPRecvBuffer = 2048

type udpAcceptor struct {
	peer.CoreSessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreCaptureIOPanic

	conn *net.UDPConn

	sesQueue *util.Queue

	sesTimeout time.Duration
}

func (self *udpAcceptor) IsReady() bool {

	return self.IsRunning()
}

func (self *udpAcceptor) Port() int {
	if self.conn == nil {
		return 0
	}

	return self.conn.LocalAddr().(*net.UDPAddr).Port
}

func (self *udpAcceptor) Start() cellnet.Peer {

	var finalAddr *util.Address
	ln, err := util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {

		addr, err := net.ResolveUDPAddr("udp", a.HostPortString(port))
		if err != nil {
			return nil, err
		}

		finalAddr = a

		return net.ListenUDP("udp", addr)
	})

	if err != nil {

		log.Errorf("#udp.resolve failed(%s) %v", self.Name(), err.Error())
		return self
	}

	self.conn = ln.(*net.UDPConn)

	if err != nil {
		log.Errorf("#udp.listen failed(%s) %s", self.Name(), err.Error())
		self.SetRunning(false)
		return self
	}

	log.Infof("#udp.listen(%s) %s", self.Name(), finalAddr.String(self.Port()))

	go self.accept()

	return self
}

func (self *udpAcceptor) protectedRecvPacket(ses *udpSession, data []byte) {
	defer func() {

		if err := recover(); err != nil {
			log.Errorf("IO panic: %s", err)
			self.conn.Close()
		}

	}()

	ses.Recv(data)
}

func (self *udpAcceptor) accept() {

	self.SetRunning(true)

	recvBuff := make([]byte, MaxUDPRecvBuffer)

	for {

		n, remoteAddr, err := self.conn.ReadFromUDP(recvBuff)
		if err != nil {
			break
		}

		if n > 0 {

			ses := self.allocSession(remoteAddr)

			if self.CaptureIOPanic() {
				self.protectedRecvPacket(ses, recvBuff[:n])
			} else {
				ses.Recv(recvBuff[:n])
			}

		}

	}

	self.SetRunning(false)

}

func (self *udpAcceptor) allocSession(addr *net.UDPAddr) *udpSession {

	var ses *udpSession

	if self.sesQueue.Count() > 0 {
		ses = self.sesQueue.Peek().(*udpSession)

		// 这个session还能用，需要重新new
		if ses.IsAlive() {
			ses = nil
		} else {
			// 可以复用
			ses = self.sesQueue.Dequeue().(*udpSession)
		}

	}

	if ses == nil {
		ses = &udpSession{}
		self.sesQueue.Enqueue(ses)
	}

	ses.timeOutTick = time.Now().Add(self.sesTimeout)
	ses.conn = self.conn
	ses.remote = addr
	ses.pInterface = self
	ses.CoreProcBundle = &self.CoreProcBundle

	return ses
}

func (self *udpAcceptor) SetSessionTTL(dur time.Duration) {
	self.sesTimeout = dur
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
		p := &udpAcceptor{
			sesQueue:   util.NewQueue(64),
			sesTimeout: time.Second,
		}

		return p
	})
}
