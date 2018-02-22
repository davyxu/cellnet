package kcp

import (
	"github.com/davyxu/cellnet"
)

type kcpPeer struct {
	cellnet.EventQueue
	cellnet.SessionManager
	*cellnet.PeerProfileImplement
	*cellnet.HandlerChainManagerImplement
}

func (k *kcpPeer) Queue() cellnet.EventQueue {
	return k.EventQueue
}

func newKcpPeer(queue cellnet.EventQueue, sm cellnet.SessionManager) *kcpPeer {
	self := &kcpPeer{
		EventQueue:						queue,
		SessionManager:					sm,
		PeerProfileImplement:			cellnet.NewPeerProfile(),
		HandlerChainManagerImplement:	cellnet.NewHandlerChainManager(),
	}
	// 设置默认发送链
	self.SetChainSend(
		cellnet.NewHandlerChain(
			cellnet.StaticEncodePacketHandler(),
		),
	)
	//设置默认读写链
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