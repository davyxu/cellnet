package socket

import (
	"net"

	"github.com/davyxu/cellnet"
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

		log.Errorln("listen failed", err.Error())
		return self
	}

	log.Debugln("listening: ", address)

	self.running = true

	// 接受线程
	go func() {
		for self.running {
			conn, err := ln.Accept()

			if err != nil {
				log.Errorln(err)
				break
			}

			ses := newSession(NewPacketStream(conn), self.EventQueue, self)

			// 添加到管理器
			self.sessionMgr.Add(ses)

			// 断开后从管理器移除
			ses.OnClose = func() {
				self.sessionMgr.Remove(ses)
			}

			// 通知逻辑
			self.PostData(NewSessionEvent(Event_SessionAccepted, ses, nil))

		}

	}()

	self.PostData(NewPeerEvent(Event_PeerStart, self))

	return self
}

func (self *socketAcceptor) Stop() {

	if !self.running {
		return
	}

	self.PostData(NewPeerEvent(Event_PeerStop, self))

	self.running = false

	self.listener.Close()
}

func NewAcceptor(pipe cellnet.EventPipe) cellnet.Peer {

	self := &socketAcceptor{
		sessionMgr:  newSessionManager(),
		peerProfile: newPeerProfile(pipe.AddQueue()),
	}

	self.PostData(NewPeerEvent(Event_PeerInit, self))

	return self
}
