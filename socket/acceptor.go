package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

type socketAcceptor struct {
	*peerProfile
	*sessionMgr

	listener net.Listener

	running bool
}

func (self *socketAcceptor) Start(address string) cellnet.Peer {

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

			// 添加到管理器
			self.sessionMgr.Add(ses)

			// 断开后从管理器移除
			ses.OnClose = func() {
				self.sessionMgr.Remove(ses)
			}

			// 通知逻辑
			self.queue.Post(NewDataEvent(Event_Accepted, ses, nil))

		}

	}()

	return self
}

func (self *socketAcceptor) Stop() {
	self.running = false

	self.listener.Close()
}

func NewAcceptor(queue *cellnet.EvQueue) cellnet.Peer {
	return &socketAcceptor{
		sessionMgr:  newSessionManager(),
		peerProfile: &peerProfile{queue: queue},
	}
}
