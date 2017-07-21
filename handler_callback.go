package cellnet

import (
	"fmt"
)

type CallbackHandler struct {
	userCallback func(*Event)
}

func (self *CallbackHandler) Call(ev *Event) {

	self.userCallback(ev)

}

func NewCallbackHandler(userCallback func(*Event)) EventHandler {

	return &CallbackHandler{
		userCallback: userCallback,
	}
}

type RegisterMessageContext struct {
	*MessageMeta
}

// 注册消息处理回调
func RegisterMessage(p Peer, msgName string, userCallback func(*Event)) *RegisterMessageContext {

	return RegisterHandler(p, msgName, NewCallbackHandler(userCallback))
}

// 注册消息处理的一系列Handler, 当有队列时, 投放到队列
func RegisterHandler(p Peer, msgName string, handlers ...EventHandler) *RegisterMessageContext {

	if p == nil {
		return nil
	}

	meta := MessageMetaByName(msgName)

	if meta == nil {
		panic(fmt.Sprintf("message register failed, name not found: %s", msgName))
	}

	if p.Queue() != nil {

		p.AddChainRecv(NewHandlerChain(
			NewMatchMsgIDHandler(meta.ID),
			StaticDecodePacketHandler(),
			NewQueuePostHandler(p.Queue(), handlers...),
		))
	} else {

		p.AddChainRecv(NewHandlerChain(
			NewMatchMsgIDHandler(meta.ID),
			StaticDecodePacketHandler(),
			handlers,
		))
	}

	return &RegisterMessageContext{MessageMeta: meta}
}

// 直接注册回调
func RegisterRawHandler(p Peer, msgName string, handlers ...EventHandler) *RegisterMessageContext {

	if p == nil {
		return nil
	}

	meta := MessageMetaByName(msgName)

	if meta == nil {
		panic(fmt.Sprintf("message register failed, %s", msgName))
	}

	p.AddChainRecv(NewHandlerChain(
		NewMatchMsgIDHandler(meta.ID),
		StaticDecodePacketHandler(),
		handlers,
	))

	return &RegisterMessageContext{MessageMeta: meta}
}
