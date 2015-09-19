package ltvsocket

import (
	"log"
	"net"

	"github.com/davyxu/cellnet"
)

type ltvAcceptor struct {
	*PeerProfile

	listener net.Listener

	running bool
}

func (self *ltvAcceptor) Start(address string) {

	ln, err := net.Listen("tcp", address)

	self.listener = ln

	if err != nil {

		log.Println("[socket] listen failed", err.Error())
		return
	}

	self.running = true

	// 接受线程
	go func() {
		for self.running {
			conn, err := ln.Accept()

			if err != nil {
				continue
			}

			ses := newSession(NewPacketStream(conn), self.queue)

			self.queue.Post(NewDataEvent(Event_Accepted, ses, nil))

		}

	}()
}

func (self *ltvAcceptor) Stop() {
	self.running = false

	self.listener.Close()
}

func init() {

	cellnet.RegisterPeerType("ltvAcceptor", func(pf *PeerProfile) Peer {
		return &ltvAcceptor{PeerProfile: pf}
	})

}
