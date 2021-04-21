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
