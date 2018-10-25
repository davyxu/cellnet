package relay

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	Ses cellnet.Session
	ack *RelayACK
	Msg interface{}
}

func (self *RecvMsgEvent) PassThroughAsInt64() int64 {
	if self.ack == nil {
		return 0
	}

	return self.ack.Int64
}

func (self *RecvMsgEvent) PassThroughAsInt64Slice() []int64 {
	if self.ack == nil {
		return nil
	}

	return self.ack.Int64Slice
}

func (self *RecvMsgEvent) PassThroughAsString() string {
	if self.ack == nil {
		return ""
	}

	return self.ack.Str
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

// 消息原路返回
func (self *RecvMsgEvent) Reply(msg interface{}) {

	// 没填的值不会被发送
	Relay(self.Ses, msg, self.ack.Int64, self.ack.Int64Slice, self.ack.Str)

}
