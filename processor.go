package cellnet

// 事件
type Event interface {

	// 事件对应的会话
	Session() Session

	// 如果消息尚未解析, 调用时将自动解析
	Message() interface{}

	// 消息ID
	MessageID() int

	// 原始数据
	MessageData() []byte
}

// 消息收发器
type MessageTransmitter interface {

	// 接收消息
	OnRecvMessage(ses Session) (ev Event, err error)

	// 发送消息
	OnSendMessage(ev Event) error
}

// 处理钩子(参数输入, 返回输出, 不给MessageProccessor处理时，可以将Event设置为nil)
type EventHooker interface {

	// 入站(接收)的事件处理
	OnInboundEvent(input Event) (output Event)

	// 出站(发送)的事件处理
	OnOutboundEvent(input Event) (output Event)
}

// 用户端处理
type EventCallback func(ev Event)
