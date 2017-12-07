package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/tcppkt"
	"net"
)

// 接受器
type socketAcceptor struct {
	internal.PeerShare

	// 保存侦听器
	listener net.Listener
}

// 异步开始侦听
func (self *socketAcceptor) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := net.Listen("tcp", self.Address())

	self.listener = ln

	if err != nil {

		log.Errorf("#listen failed(%s) %v", self.NameOrAddress(), err.Error())
		return self
	}

	log.Infof("#listen(%s) %s", self.Name(), self.Address())

	go self.accept()

	return self
}

func (self *socketAcceptor) accept() {
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

func (self *socketAcceptor) onNewSession(conn net.Conn) {

	ses := internal.NewSession(conn, &self.PeerShare, nil)

	ses.(interface {
		Start()
	}).Start()

	self.FireEvent(cellnet.SessionAcceptedEvent{ses})
}

func (self *socketAcceptor) IsAcceptor() bool {
	return true
}

// 停止侦听器
func (self *socketAcceptor) Stop() {
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

	cellnet.RegisterPeerCreator("ltv.tcp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &socketAcceptor{}

		initEvent(&config)

		p.Init(p, config)

		return p
	})
}
