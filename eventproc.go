package cellnet

type Event interface {
	Session() Session
	Message() interface{}
}

// 消息处理器
type MessageProcessor interface {
	OnRecvMessage(ses Session) (raw interface{}, err error)
	OnSendMessage(ses Session, raw interface{}) error
}

// 处理钩子(参数输入, 返回输出, 不给MessageProccessor处理时，可以将Event设置为nil)
type EventHooker interface {
	OnInboundEvent(input Event) (output Event)
	OnOutboundEvent(input Event) (output Event)
}

// 用户端处理
type EventHandler interface {
	OnEvent(ev Event)
}

// 直接回调用户回调
type UserMessageHandler func(ev Event)

func (self UserMessageHandler) OnEvent(ev Event) {

	if self != nil {
		self(ev)
	}
}

// 放队列中回调
type UserMessageHandlerQueued func(ev Event)

func (self UserMessageHandlerQueued) OnEvent(ev Event) {

	if self != nil {
		SessionQueuedCall(ev.Session(), func() {

			self(ev)
		})
	}

}
