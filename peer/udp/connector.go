package udp

import (
	cellevent "github.com/davyxu/cellnet/event"
	"net"
	"sync/atomic"
)

type Connector struct {
	*Peer

	Address string

	remoteAddr *net.UDPAddr

	// 连接会话
	Session *Session

	closed int64
}

func (self *Connector) Connect(address string) error {
	self.Address = address

	var err error
	self.remoteAddr, err = net.ResolveUDPAddr("udp", address)

	if err != nil {
		return err
	}

	return self.conn()
}

func (self *Connector) AsyncConnect(address string) error {
	self.Address = address
	var err error
	self.remoteAddr, err = net.ResolveUDPAddr("udp", address)

	if err != nil {
		return err
	}

	go self.conn()

	return nil
}

func (self *Connector) conn() error {
	conn, err := net.DialUDP("udp", nil, self.remoteAddr)
	if err != nil {
		self.ProcEvent(cellevent.BuildSystemEvent(self.Session, &cellevent.SessionConnectError{}))
		return err
	}

	atomic.StoreInt64(&self.closed, 0)

	self.Session.conn = conn

	ses := self.Session

	self.ProcEvent(cellevent.BuildSystemEvent(self.Session, &cellevent.SessionConnected{}))

	recvBuff := make([]byte, MaxUDPRecvBuffer)

	for {

		if atomic.LoadInt64(&self.closed) != 0 {
			break
		}

		n, _, err := conn.ReadFromUDP(recvBuff)
		if err != nil {
			break
		}

		if n > 0 {
			ses.Recv(recvBuff[:n])
		}

	}

	if self.Session.conn != nil {
		self.Session.conn.Close()
	}

	self.ProcEvent(cellevent.BuildSystemEvent(self.Session, &cellevent.SessionClosed{}))
	return nil
}

func (self *Connector) Close() {
	atomic.StoreInt64(&self.closed, 1)
}

func (self *Connector) Port() int {
	if self.Session.conn == nil {
		return 0
	}
	return self.Session.conn.LocalAddr().(*net.UDPAddr).Port
}

func NewConnector() *Connector {
	self := &Connector{
		Peer: newPeer(),
	}

	ses := &Session{}
	ses.parent = self
	ses.peer = self.Peer

	self.Session = ses

	return self
}
