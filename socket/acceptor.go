package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

type ltvAcceptor struct {
	*peerProfile
	*sessionMgr

	listener net.Listener

	running bool
}

func (self *ltvAcceptor) Start(address string) cellnet.Peer {

	ln, err := net.Listen("tcp", address)

	self.listener = ln

	if err != nil {

		log.Println("[socket] listen failed", err.Error())
		return self
	}

	log.Println("listening: ", address)

	self.running = true

	// 接受线程
	go func() {
		for self.running {
			conn, err := ln.Accept()

			if err != nil {
				continue
			}

			ses := newSession(NewPacketStream(conn), self.queue, self)

			self.sessionMgr.Add(ses)

			ses.OnClose = func() {
				self.sessionMgr.Remove(ses)
			}

			self.queue.Post(NewDataEvent(Event_Accepted, ses, nil))

		}

	}()

	return self
}

func (self *ltvAcceptor) Stop() {
	self.running = false

	self.listener.Close()
}

func NewAcceptor(queue *cellnet.EvQueue) cellnet.Peer {
	return &ltvAcceptor{
		sessionMgr:  newSessionManager(),
		peerProfile: &peerProfile{queue: queue},
	}
}
