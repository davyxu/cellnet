package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"net"
	"strings"
)

// 接受器
type tcpAcceptor struct {
	peer.SessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	// 保存侦听器
	listener net.Listener
}

func (self *tcpAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *tcpAcceptor) IsReady() bool {

	return self.IsRunning()
}

// 异步开始侦听
func (self *tcpAcceptor) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := util.DetectPort(self.Address(), func(s string) (interface{}, error) {
		return net.Listen("tcp", s)
	})

	if err != nil {

		log.Errorf("#tcp.listen failed(%s) %v", self.NameOrAddress(), err.Error())

		self.SetRunning(false)

		return self
	}

	self.listener = ln.(net.Listener)

	log.Infof("#tcp.listen(%s) %s", self.Name(), self.ListenAddress())

	go self.accept()

	return self
}

func (self *tcpAcceptor) ListenAddress() string {

	pos := strings.Index(self.Address(), ":")
	if pos == -1 {
		return self.Address()
	}

	host := self.Address()[:pos]

	return util.JoinAddress(host, self.Port())
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
				log.Errorf("#tcp.accept failed(%s) %v", self.NameOrAddress(), err.Error())
			}

			break
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onNewSession(conn)

	}

	self.SetRunning(false)

	self.EndStopping()

}

func (self *tcpAcceptor) onNewSession(conn net.Conn) {

	self.ApplySocketOption(conn)

	ses := newSession(conn, self, nil)

	ses.Start()

	self.PostEvent(&cellnet.RecvMsgEvent{ses, &cellnet.SessionAccepted{}})
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

func (self *tcpAcceptor) TypeName() string {
	return "tcp.Acceptor"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &tcpAcceptor{
			SessionManager: new(peer.CoreSessionManager),
		}

		p.CoreTCPSocketOption.Init()

		return p
	})
}
