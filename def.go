package cellnet

import "github.com/davyxu/x/container"

// 会话
type Session interface {
	// 发送消息
	Send(msg any)
}

// 事件
type Event interface {

	// 事件对应的会话
	Session() Session

	// 如果消息尚未解析, 调用时将自动解析
	Message() any

	// 消息ID
	MessageId() int

	// 原始数据
	MessageData() []byte
}

type Codec interface {
	// 将数据转换为字节数组
	Encode(msgObj any, ps *xcontainer.Mapper) (data any, err error)

	// 将字节数组转换为数据
	Decode(data any, msgObj any) error

	Name() string
}
