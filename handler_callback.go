package cellnet

func MessageRegistedCount(evd EventDispatcher, msgName string) int {

	msgMeta := MessageMetaByName(msgName)
	if msgMeta == nil {
		return 0
	}

	return evd.CountByID(int(msgMeta.ID))
}

type RegisterMessageContext struct {
	*MessageMeta
	*HandlerContext
}

type CallbackHandler struct {
	BaseEventHandler
	userCallback func(*SessionEvent)
}

func (self *CallbackHandler) Call(ev *SessionEvent) error {

	self.userCallback(ev)

	return self.CallNext(ev)
}

func NewCallbackHandler(userCallback func(*SessionEvent)) EventHandler {
	return &CallbackHandler{
		userCallback: userCallback,
	}
}

// 注册消息处理回调
// DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func RegisterMessage(dh EventDispatcher, msgName string, userCallback func(*SessionEvent)) *RegisterMessageContext {

	meta := MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("message register failed, %s", msgName)
		return nil
	}

	ctx := dh.AddHandler(int(meta.ID), LinkHandler(NewDecodePacketHandler(meta), NewCallbackHandler(userCallback)))

	return &RegisterMessageContext{MessageMeta: meta, HandlerContext: ctx}
}

// 注册消息处理的一系列Handler
// DispatcherHandler -> socket.DecodePacketHandler -> ...
func RegisterHandler(dh EventDispatcher, msgName string, handlers ...EventHandler) *RegisterMessageContext {

	meta := MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("message register failed, %s", msgName)
		return nil
	}

	decoder := NewDecodePacketHandler(meta)

	if len(handlers) > 0 {
		decoder.SetNext(LinkHandler(handlers...))
	}

	ctx := dh.AddHandler(int(meta.ID), decoder)

	return &RegisterMessageContext{MessageMeta: meta, HandlerContext: ctx}
}
