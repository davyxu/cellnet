package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

type ltvAcceptor struct {
	listener net.Listener

	running bool

	queue *cellnet.EvQueue
}

func (self *ltvAcceptor) Start(address string) {

	ln, err := net.Listen("tcp", address)

	self.listener = ln

	if err != nil {

		log.Println("[socket] listen failed", err.Error())
		return
	}

	self.running = true

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

func newAcceptor(evq *cellnet.EvQueue) *ltvAcceptor {
	return &ltvAcceptor{queue: evq}
}
