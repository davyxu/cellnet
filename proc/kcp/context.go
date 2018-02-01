package kcp

import (
	"github.com/davyxu/cellnet"
	"time"
)

type kcpContext struct {
	kcp *KCP

	// kcp receiving is based on packets
	// recvbuf turns packets into stream
	recvbuf []byte
	bufptr  []byte

	readSignal     chan struct{}
	exitTickSignal chan struct{}

	userFunc cellnet.EventProc

	ses cellnet.Session

	closed bool
}

func (self *kcpContext) Close() {
	close(self.readSignal)
	self.closed = true
	self.exitTickSignal <- struct{}{}
}

func (self *kcpContext) tickLoop() {

	ticker := time.NewTicker(time.Millisecond * 10)
	for {

		select {
		case <-ticker.C:
			self.kcp.Update()

		case <-self.exitTickSignal:
			return
		}

	}

}

func newContext(ses cellnet.Session, userFunc cellnet.EventProc) *kcpContext {

	var self *kcpContext

	self = &kcpContext{
		userFunc:       userFunc,
		ses:            ses,
		recvbuf:        make([]byte, mtuLimit),
		readSignal:     make(chan struct{}, 1),
		exitTickSignal: make(chan struct{}),
		kcp: NewKCP(0, func(buf []byte, size int) {

			if size >= IKCP_OVERHEAD {
				self.output(buf[:size])
			}
		}),
	}

	go self.recvLoop()

	go self.tickLoop()

	return self
}
