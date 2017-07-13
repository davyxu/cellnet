package socket

import (
	"net"

	"github.com/davyxu/cellnet"
)

type socketAcceptor struct {
	*peerBase

	listener net.Listener
}

func (self *socketAcceptor) Start(address string) cellnet.Peer {

	self.waitStopFinished()

	if self.IsRunning() {
		return self
	}

	self.address = address

	ln, err := net.Listen("tcp", address)

	self.listener = ln

	if err != nil {

		log.Errorf("#listen failed(%s) %v", self.nameOrAddress(), err.Error())
		return self
	}

	log.Infof("#listen(%s) %s", self.name, self.address)

	// 接受线程
	go self.accept()

	return self
}

func (self *socketAcceptor) accept() {

	self.SetRunning(true)

	for {
		conn, err := self.listener.Accept()

		if self.isStopping() {
			break
		}

		if err != nil {

			// 调试状态时, 才打出accept的具体错误
			if log.IsDebugEnabled() {
				log.Errorf("#accept failed(%s) %v", self.nameOrAddress(), err.Error())
			}

			systemError(nil, cellnet.Event_AcceptFailed, errToResult(err), self.safeRecvHandler())

			break
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go func() {

			ses := newSession(self.genPacketStream(conn), self)

			// 添加到管理器
			self.SessionManager.Add(ses)

			// 断开后从管理器移除
			ses.OnClose = func() {
				self.SessionManager.Remove(ses)
			}

			ses.run()

			// 通知逻辑
			systemEvent(ses, cellnet.Event_Accepted, self.safeRecvHandler())
		}()

	}

	self.SetRunning(false)

	self.endStopping()
}

func (self *socketAcceptor) Stop() {

	if !self.IsRunning() {
		return
	}

	if self.isStopping() {
		return
	}

	self.startStopping()

	self.listener.Close()

	// 断开所有连接
	self.CloseAllSession()

	// 等待线程结束
	self.waitStopFinished()
}

func NewAcceptor(q cellnet.EventQueue) cellnet.Peer {

	self := &socketAcceptor{
		peerBase: newPeerBase(q, NewSessionManager()),
	}

	return self
}
