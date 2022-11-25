package cellevent

import (
	"github.com/davyxu/cellnet"
)

// 接收到消息
type RecvMsg struct {
	Ses cellnet.Session
	Msg any

	// 原始数据
	MsgID   int
	MsgData []byte
}

func (self *RecvMsg) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsg) MessageID() int {
	return self.MsgID
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
