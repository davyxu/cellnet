package cellnet

import (
	"fmt"
)

type CallbackHandler struct {
	userCallback func(*SessionEvent)
}

func (self *CallbackHandler) Call(ev *SessionEvent) {

	self.userCallback(ev)

}

func NewCallbackHandler(userCallback func(*SessionEvent)) EventHandler {

	return &CallbackHandler{
		userCallback: userCallback,
	}
}

type RegisterMessageContext struct {
	*MessageMeta
}

// 注册消息处理回调
// DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func RegisterMessage(p Peer, msgName string, userCallback func(*SessionEvent)) *RegisterMessageContext {

	return RegisterHandler(p, msgName, NewCallbackHandler(userCallback))
}

// 注册消息处理的一系列Handler
// DispatcherHandler -> socket.DecodePacketHandler -> ...
func RegisterHandler(p Peer, msgName string, handlers ...EventHandler) *RegisterMessageContext {

	if p == nil {
		return nil
	}

	meta := MessageMetaByName(msgName)

	if meta == nil {
		panic(fmt.Sprintf("message register failed, %s", msgName))
	}

	if p.Queue() != nil {
		p.AddHandler(int(meta.ID), HandlerLink(NewQueuePostHandler(p.Queue(), handlers)))
	} else {
		p.AddHandler(int(meta.ID), HandlerLink(handlers))
	}

	return &RegisterMessageContext{MessageMeta: meta}
}
