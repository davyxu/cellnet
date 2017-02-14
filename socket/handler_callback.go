package socket

import "github.com/davyxu/cellnet"

func MessageRegistedCount(evd cellnet.EventDispatcher, msgName string) int {

	msgMeta := cellnet.MessageMetaByName(msgName)
	if msgMeta == nil {
		return 0
	}

	return evd.CountByID(int(msgMeta.ID))
}

type RegisterMessageContext struct {
	*cellnet.MessageMeta
	*cellnet.HandlerContext
}

type CallbackHandler struct {
	cellnet.BaseEventHandler
	userCallback func(*cellnet.SessionEvent)
}

func (self *CallbackHandler) Call(ev *cellnet.SessionEvent) error {

	self.userCallback(ev)

	return self.CallNext(ev)
}

func NewCallbackHandler(userCallback func(*cellnet.SessionEvent)) cellnet.EventHandler {
	return &CallbackHandler{
		userCallback: userCallback,
	}
}

func RegisterMessage(dh cellnet.EventDispatcher, msgName string, userCallback func(*cellnet.SessionEvent)) *RegisterMessageContext {

	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("message register failed, %s", msgName)
		return nil
	}

	ctx := dh.AddHandler(int(meta.ID), cellnet.LinkHandler(NewDecodePacketHandler(meta), NewCallbackHandler(userCallback)))

	return &RegisterMessageContext{MessageMeta: meta, HandlerContext: ctx}
}
