package cellevent

import (
	"github.com/davyxu/cellnet"
)

// 接收到消息
type RecvMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}

	// 原始数据
	MsgID   int
	MsgData []byte
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsgEvent) MessageID() int {
	return self.MsgID
}

func (self *RecvMsgEvent) MessageData() []byte {
	return self.MsgData
}

func (self *RecvMsgEvent) Message() interface{} {

	if self.Msg == nil {
		self.Msg = InternalDecodeHandler(self)
	}

	return self.Msg
}

// 会话开始发送数据事件
type SendMsgEvent struct {
	Ses cellnet.Session
	Msg interface{} // 用户需要发送的消息

	// 原始数据
	MsgID   int
	MsgData []byte
}

func (self *SendMsgEvent) Message() interface{} {
	if self.Msg == nil {
		self.Msg = InternalDecodeHandler(self)
	}

	return self.Msg
}

func (self *SendMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *SendMsgEvent) MessageID() int {
	return self.MsgID
}

func (self *SendMsgEvent) MessageData() []byte {
	return self.MsgData
}

// rpc, relay, 普通消息
type ReplyEvent interface {
	Reply(msg interface{})
}

var InternalDecodeHandler func(ev cellnet.Event) (msg interface{})
