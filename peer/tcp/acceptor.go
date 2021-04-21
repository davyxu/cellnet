package tcp

import (
	"errors"
	cellevent "github.com/davyxu/cellnet/event"
	xnet "github.com/davyxu/x/net"
	"net"
	"time"
)

var (
	ErrClosed = errors.New("acceptor closed")
)

type Acceptor struct {
	*Peer

	// 连接地址
	Address string

	listener net.Listener

	doneChan chan struct{}
}

func (self *Acceptor) Listen(addr string) error {
	self.Address = addr
	ln, err := xnet.DetectPort(self.Address, func(a *xnet.Address, port int) (interface{}, error) {
		return net.Listen("tcp", a.HostPortString(port))
	})

	if err != nil {
		return err
	}

	self.listener = ln.(net.Listener)

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
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *Acceptor) Accept() error {

	if self.listener == nil {
		return ErrClosed
	}

	var tempDelay time.Duration

	for {
		conn, err := self.listener.Accept()

		if err != nil {

			select {
			case <-self.doneChan:
				return ErrClosed
			default:
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
				continue
			}

			return err
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onNewSession(conn)
	}
}

func (self *Acceptor) onNewSession(conn net.Conn) {

	self.ApplySocketOption(conn)

	ses := newSession(conn, self.Peer, self)

	ses.Start()

	self.ProcEvent(cellevent.BuildSystemEvent(ses, &cellevent.SessionAccepted{}))
}

func (self *Acceptor) Close() error {
	if self.listener == nil {
		return nil
	}

	self.closeDone()

	err := self.listener.Close()

	self.CloseAll()

	return err
}

func (self *Acceptor) closeDone() {
	select {
	case <-self.doneChan:
		// 已经关闭, 不要重复关闭
	default:
		close(self.doneChan)
	}
}

func NewAcceptor() *Acceptor {
	self := &Acceptor{
		Peer:     newPeer(),
		doneChan: make(chan struct{}),
	}

	self.Init()

	return self
}
