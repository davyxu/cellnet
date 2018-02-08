package cellnet

// 接收到消息
type RecvMsgEvent struct {
	Ses BaseSession
	Msg interface{}
}

func (self *RecvMsgEvent) BaseSession() BaseSession {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Send(msg interface{}) {
	self.Ses.Send(msg)
}

// 会话开始发送数据事件
type SendMsgEvent struct {
	Ses BaseSession
	Msg interface{} // 用户需要发送的消息
}

func (self *SendMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *SendMsgEvent) BaseSession() BaseSession {
	return self.Ses
}
