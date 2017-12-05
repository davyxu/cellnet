package tcppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/tcppkt"
	"net"
)

type socketConnector struct {
	internal.PeerShare

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

	ses := internal.NewSession(conn, &self.PeerShare)
	self.ses = ses

	// 发生错误时退出
	if err != nil {
		self.FireEvent(cellnet.ConnectErrorEvent{ses, err})
		return
	}

	ses.(interface {
		Start()
	}).Start()

}

func init() {

	cellnet.RegisterPeerCreator("tcp.Connector", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &socketConnector{}
		config.Event = tcppkt.ProcTLVPacket(msglog.ProcMsgLog(rpc.ProcRPC(config.Event)))

		p.Init(p, config)

		return p
	})
}
