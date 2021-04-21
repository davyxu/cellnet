package udp

import (
	xnet "github.com/davyxu/x/net"
	"net"
	"time"
)

type Acceptor struct {
	*Peer

	// 连接地址
	Address string

	conn *net.UDPConn

	sesTimeout       time.Duration
	sesCleanTimeout  time.Duration
	sesCleanLastTime time.Time

	sesByConnTrack map[connTrackKey]*Session
}

func (self *Acceptor) Listen(addr string) error {
	self.Address = addr
	ln, err := xnet.DetectPort(self.Address, func(a *xnet.Address, port int) (interface{}, error) {
		addr, err := net.ResolveUDPAddr("udp", a.HostPortString(port))
		if err != nil {
			return nil, err
		}

		return net.ListenUDP("udp", addr)
	})

	if err != nil {
		return err
	}

	self.conn = ln.(*net.UDPConn)

	return nil
}
func (self *Acceptor) ListenAndAccept(addr string) error {

	err := self.Listen(addr)
	if err != nil {
		return err
	}

	go self.Accept()

	return nil
}

func (self *Acceptor) ListenPort() int {
	if self.conn == nil {
		return 0
	}

	return self.conn.LocalAddr().(*net.UDPAddr).Port
}

const MaxUDPRecvBuffer = 2048

func (self *Acceptor) Accept() error {
	recvBuff := make([]byte, MaxUDPRecvBuffer)

	for {

		n, remoteAddr, err := self.conn.ReadFromUDP(recvBuff)
		if err != nil {
			break
		}

		self.checkTimeoutSession()

		if n > 0 {

			ses := self.getSession(remoteAddr)
			ses.Recv(recvBuff[:n])
		}

	}

	return nil
}

// 检查超时session
func (self *Acceptor) checkTimeoutSession() {
	now := time.Now()

	// 定时清理超时的session
	if now.After(self.sesCleanLastTime.Add(self.sesCleanTimeout)) {
		sesToDelete := make([]*Session, 0, 10)
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

func (self *Acceptor) getSession(addr *net.UDPAddr) *Session {

	key := newConnTrackKey(addr)

	ses := self.sesByConnTrack[*key]

	if ses == nil {
		ses = &Session{}
		ses.conn = self.conn
		ses.remote = addr
		ses.parent = self
		ses.peer = self.Peer
		ses.key = key
		self.sesByConnTrack[*key] = ses
	}

	// 续租
	ses.timeOutTick = time.Now().Add(self.sesTimeout)

	return ses
}

func NewAcceptor() *Acceptor {
	self := &Acceptor{
		Peer:           newPeer(),
		sesByConnTrack: map[connTrackKey]*Session{},
	}

	return self
}
