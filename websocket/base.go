package websocket

import (
	"github.com/davyxu/cellnet"
)

// Peer间的共享数据
type wsPeer struct {
	cellnet.EventQueue

	cellnet.SessionManager

	*cellnet.PeerProfileImplement

	*cellnet.HandlerChainManagerImplement
}

func (self *wsPeer) Queue() cellnet.EventQueue {
	return self.EventQueue
}

func newPeer(queue cellnet.EventQueue, sm cellnet.SessionManager) *wsPeer {

	self := &wsPeer{
		EventQueue:                   queue,
		SessionManager:               sm,
		PeerProfileImplement:         cellnet.NewPeerProfile(),
		HandlerChainManagerImplement: cellnet.NewHandlerChainManager(),
	}

	self.SetChainSend(
		cellnet.NewHandlerChain(
			cellnet.StaticEncodePacketHandler(),
		),
	)

	return self
}
