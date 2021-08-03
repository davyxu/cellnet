package cellnet

import xframe "github.com/davyxu/x/frame"

// 会话
type Session interface {
	// 发送消息
	Send(msg interface{})
}

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

type Codec interface {
	// 将数据转换为字节数组
	Encode(msgObj interface{}, ps *xframe.Mapper) (data interface{}, err error)

	// 将字节数组转换为数据
	Decode(data interface{}, msgObj interface{}) error

	Name() string
}
