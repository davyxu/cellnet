package socket

import (
	"github.com/davyxu/cellnet"
	"net"
)

// Peer间的共享数据
type socketPeer struct {
	cellnet.EventQueue
	// 会话管理器
	cellnet.SessionManager

	// 共享配置
	*cellnet.PeerProfileImplement

	// 处理链管理
	*cellnet.HandlerChainManagerImplement

	// socket配置
	*socketOptions

	// 停止过程同步
	stopping chan bool
}

func (self *socketPeer) waitStopFinished() {
	// 如果正在停止时, 等待停止完成
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}

func (self *socketPeer) isStopping() bool {
	return self.stopping != nil
}

func (self *socketPeer) startStopping() {
	self.stopping = make(chan bool)
}

func (self *socketPeer) endStopping() {
	select {
	case self.stopping <- true:

	default:
		self.stopping = nil
	}
}

func (self *socketPeer) Queue() cellnet.EventQueue {
	return self.EventQueue
}

func newSocketPeer(queue cellnet.EventQueue, sm cellnet.SessionManager) *socketPeer {

	self := &socketPeer{
		EventQueue:                   queue,
		SessionManager:               sm,
		socketOptions:                newSocketOptions(),
		PeerProfileImplement:         cellnet.NewPeerProfile(),
		HandlerChainManagerImplement: cellnet.NewHandlerChainManager(),
	}

	// 设置默认发送链
	self.SetChainSend(
		cellnet.NewHandlerChain(
			cellnet.StaticEncodePacketHandler(),
		),
	)

	// 设置默认读写链
	self.SetReadWriteChain(func() *cellnet.HandlerChain {
		return cellnet.NewHandlerChain(
			cellnet.NewFixedLengthFrameReader(10),
			NewPrivatePacketReader(),
		)
	}, func() *cellnet.HandlerChain {
		return cellnet.NewHandlerChain(NewPrivatePacketWriter(),
			cellnet.NewFixedLengthFrameWriter(),
		)
	})

	return self
}

func errToResult(err error) cellnet.Result {

	if err == nil {
		return cellnet.Result_OK
	}

	switch n := err.(type) {
	case net.Error:
		if n.Timeout() {
			return cellnet.Result_SocketTimeout
		}
	}

	return cellnet.Result_SocketError
}
