package cellnet

type Event interface {
	Session() Session
	Message() interface{}
}

// 消息收发器
type MessageTransmitter interface {
	OnRecvMessage(ses Session) (raw interface{}, err error)
	OnSendMessage(ses Session, raw interface{}) error
}

// 处理钩子(参数输入, 返回输出, 不给MessageProccessor处理时，可以将Event设置为nil)
type EventHooker interface {
	OnInboundEvent(input Event) (output Event)
	OnOutboundEvent(input Event) (output Event)
}

// 用户端处理
type EventCallback func(ev Event)

// 直接回调用户回调

// 放队列中回调
func NewQueuedEventCallback(callback EventCallback) EventCallback {

	return func(ev Event) {
		if callback != nil {
			SessionQueuedCall(ev.Session(), func() {

				callback(ev)
			})
		}
	}

}
