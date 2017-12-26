package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
)

// 接受器
type tcpAcceptor struct {
	internal.PeerShare

	// 保存侦听器
	listener net.Listener
}

// 异步开始侦听
func (self *tcpAcceptor) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := net.Listen("tcp", self.Address())

	if err != nil {

		log.Errorf("#listen failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	self.listener = ln

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go self.accept()

	return self
}

func (self *tcpAcceptor) accept() {
	self.SetRunning(true)

	for {
		conn, err := self.listener.Accept()

		if self.IsStopping() {
			break
		}

		if err != nil {

			// 调试状态时, 才打出accept的具体错误
			if log.IsDebugEnabled() {
				log.Errorf("#accept failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			//extend.PostSystemEvent(nil, cellnet.Event_AcceptFailed, self.ChainListRecv(), errToResult(err))

			break
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onNewSession(conn)

	}

	self.SetRunning(false)

	self.EndStopping()

}

func (self *tcpAcceptor) onNewSession(conn net.Conn) {

	ses := newTCPSession(conn, &self.PeerShare, nil)

	ses.(interface {
		Start()
	}).Start()

	self.CallInboundProc(&cellnet.RecvMsgEvent{ses, &comm.SessionAccepted{}})
}

func (self *tcpAcceptor) IsAcceptor() bool {
	return true
}

// 停止侦听器
func (self *tcpAcceptor) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	self.listener.Close()

	// 断开所有连接
	self.CloseAllSession()

	// 等待线程结束
	self.WaitStopFinished()
}

func init() {

	cellnet.RegisterPeerCreator("tcp.Acceptor", func() cellnet.Peer {
		p := &tcpAcceptor{}

		p.Init(p)

		return p
	})
}
