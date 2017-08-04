package socket

import (
	"net"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
)

type socketAcceptor struct {
	*socketPeer

	listener net.Listener
}

func (self *socketAcceptor) Start(address string) cellnet.Peer {

	self.waitStopFinished()

	if self.IsRunning() {
		return self
	}

	self.SetAddress(address)

	ln, err := net.Listen("tcp", address)

	self.listener = ln

	if err != nil {

		log.Errorf("#listen failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

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
				log.Errorf("#accept failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			extend.PostSystemEvent(nil, cellnet.Event_AcceptFailed, self.ChainListRecv(), errToResult(err))

			break
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onAccepted(conn)

	}

	self.SetRunning(false)

	self.endStopping()
}

func (self *socketAcceptor) onAccepted(conn net.Conn) {

	ses := newSession(conn, self)

	// 添加到管理器
	self.Add(ses)

	// 断开后从管理器移除
	ses.OnClose = func() {
		self.Remove(ses)
	}

	ses.run()

	// 通知逻辑
	extend.PostSystemEvent(ses, cellnet.Event_Accepted, self.ChainListRecv(), cellnet.Result_OK)
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
		socketPeer: newSocketPeer(q, cellnet.NewSessionManager()),
	}

	return self
}
