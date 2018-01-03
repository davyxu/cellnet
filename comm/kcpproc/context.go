package kcpproc

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

	readEvent chan struct{}

	userFunc cellnet.EventProc

	ses cellnet.Session
}



func (self *kcpContext) tickLoop() {

	heatbeat := time.NewTicker(time.Millisecond * 10)
	for {

		select {
		case <-heatbeat.C:

			self.kcp.Update()
		}

	}

}

func newContext(ses cellnet.Session, userFunc cellnet.EventProc) *kcpContext {

	var self *kcpContext

	self = &kcpContext{
		userFunc:  userFunc,
		ses:       ses,
		recvbuf : make([]byte, mtuLimit),
		readEvent: make(chan struct{}, 1),
		kcp: NewKCP(0, func(buf []byte, size int) {

			if size >= IKCP_OVERHEAD {
				self.output(buf[:size])
			}
		}),
	}

	self.kcp.WndSize(128, 128)
	self.kcp.NoDelay(0, 10, 0, 0)

	go self.recvLoop()

	go self.tickLoop()

	return self
}
