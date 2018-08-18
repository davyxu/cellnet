package relay

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}

	PassThrough interface{}
}

func (self *RecvMsgEvent) PassThroughAsInt64() int64 {

	if v, ok := self.PassThrough.(int64); ok {
		return v
	}

	return 0
}

func (self *RecvMsgEvent) PassThroughAsInt64Slice() []int64 {

	if v, ok := self.PassThrough.([]int64); ok {
		return v
	}

	return nil
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

// 消息原路返回
func (self *RecvMsgEvent) Reply(msg interface{}) {

	Relay(self.Ses, msg, self.PassThrough)
}
