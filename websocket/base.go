package websocket

import (
	"github.com/davyxu/cellnet"
)

// Peer间的共享数据
type wsPeer struct {
	cellnet.EventQueue
	// 会话管理器
	cellnet.SessionManager

	// 共享配置
	*cellnet.BasePeerImplement

	// 自带派发器
	*cellnet.DispatcherHandler
}

func (self *wsPeer) Queue() cellnet.EventQueue {
	return self.EventQueue
}

func newPeer(queue cellnet.EventQueue, sm cellnet.SessionManager) *wsPeer {

	self := &wsPeer{
		EventQueue:        queue,
		DispatcherHandler: cellnet.NewDispatcherHandler(),
		SessionManager:    sm,
		BasePeerImplement: cellnet.NewBasePeer(),
	}

	self.BasePeerImplement.SetHandlerList(BuildRecvHandler(self.DispatcherHandler), BuildSendHandler())

	return self
}
