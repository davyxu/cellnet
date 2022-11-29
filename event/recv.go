package cellevent

import (
	"github.com/davyxu/cellnet"
)

// 接收到消息
type RecvMsg struct {
	Ses cellnet.Session
	Msg any

	// 原始数据
	MsgId   int
	MsgData []byte
}

func (self *RecvMsg) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsg) MessageId() int {
	return self.MsgId
}

func (self *RecvMsg) MessageData() []byte {
	return self.MsgData
}

func (self *RecvMsg) Message() any {

	if self.Msg == nil {
		self.Msg = InternalDecodeHandler(self)
	}

	return self.Msg
}

func (self *RecvMsg) Send(msg any) {
	if self.Ses != nil {
		self.Ses.Send(msg)
	}
}
