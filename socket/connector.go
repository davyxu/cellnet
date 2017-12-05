package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
)

type socketConnector struct {
	socketPeer
	internal.SessionManager

	ses cellnet.Session
}

func (self *socketConnector) Start() cellnet.Peer {

	go self.connect(self.PeerAddress)

	return self
}

func (self *socketConnector) Session() cellnet.Session {
	return self.ses
}

func (self *socketConnector) Stop() {

}

// 连接器，传入连接地址和发送封包次数
func (self *socketConnector) connect(address string) {

	// 尝试用Socket连接地址
	conn, err := net.Dial("tcp", address)

	ses := newSession(conn, &self.socketPeer)
	self.ses = ses

	// 发生错误时退出
	if err != nil {
		self.fireEvent(ConnectErrorEvent{ses, err})
		return
	}

	ses.start()

}

func init() {

	cellnet.RegisterPeerCreator("tcp.Connector", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &socketConnector{
			SessionManager: internal.NewSessionManager(),
		}

		p.PeerConfig = config
		p.peerInterface = p

		return p
	})
}
