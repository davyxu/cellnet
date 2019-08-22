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

	sesTimeout       time.Duration
	sesCleanTimeout  time.Duration
	sesCleanLastTime time.Time

	sesByConnTrack map[connTrackKey]*udpSession
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

		log.Errorf("#udp.listen failed(%s) %v", self.Name(), err.Error())
		return self
	}

	self.conn = ln.(*net.UDPConn)

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

		self.checkTimeoutSession()

		if n > 0 {

			ses := self.getSession(remoteAddr)

			if self.CaptureIOPanic() {
				self.protectedRecvPacket(ses, recvBuff[:n])
			} else {
				ses.Recv(recvBuff[:n])
			}

		}

	}

	self.SetRunning(false)

}

// 检查超时session
func (self *udpAcceptor) checkTimeoutSession() {
	now := time.Now()

	// 定时清理超时的session
	if now.After(self.sesCleanLastTime.Add(self.sesCleanTimeout)) {
		sesToDelete := make([]*udpSession, 0, 10)
		for _, ses := range self.sesByConnTrack {
			if !ses.IsAlive() {
				sesToDelete = append(sesToDelete, ses)
			}
		}

		for _, ses := range sesToDelete {
			delete(self.sesByConnTrack, *ses.key)
		}

		self.sesCleanLastTime = now
	}
}

func (self *udpAcceptor) getSession(addr *net.UDPAddr) *udpSession {

	key := newConnTrackKey(addr)

	ses := self.sesByConnTrack[*key]

	if ses == nil {
		ses = &udpSession{}
		ses.conn = self.conn
		ses.remote = addr
		ses.pInterface = self
		ses.CoreProcBundle = &self.CoreProcBundle
		ses.key = key
		self.sesByConnTrack[*key] = ses
	}

	// 续租
	ses.timeOutTick = time.Now().Add(self.sesTimeout)

	return ses
}

func (self *udpAcceptor) SetSessionTTL(dur time.Duration) {
	self.sesTimeout = dur
}

func (self *udpAcceptor) SetSessionCleanTimeout(dur time.Duration) {
	self.sesCleanTimeout = dur
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
			sesTimeout:       time.Minute,
			sesCleanTimeout:  time.Minute,
			sesCleanLastTime: time.Now(),
			sesByConnTrack:   make(map[connTrackKey]*udpSession),
		}

		return p
	})
}
