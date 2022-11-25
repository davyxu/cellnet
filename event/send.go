package cellevent

import "github.com/davyxu/cellnet"

// 会话开始发送数据事件
type SendMsg struct {
	Ses cellnet.Session
	Msg any // 用户需要发送的消息

	// 原始数据
	MsgID   int
	MsgData []byte
}

func (self *SendMsg) Message() any {
	if self.Msg == nil {
		self.Msg = InternalDecodeHandler(self)
	}

	return self.Msg
}

func (self *SendMsg) Session() cellnet.Session {
	return self.Ses
}

func (self *SendMsg) MessageID() int {
	return self.MsgID
}

func (self *SendMsg) MessageData() []byte {
	return self.MsgData
}
