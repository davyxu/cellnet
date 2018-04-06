package relay

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}

	ContextID []int64
}

func (self *RecvMsgEvent) OneContextID() int64 {
	if len(self.ContextID) == 0 {
		return 0
	}

	return self.ContextID[0]
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

// 消息原路返回
func (self *RecvMsgEvent) RelayBack(msg interface{}) {

	Relay(self.Ses, msg, self.ContextID...)
}
