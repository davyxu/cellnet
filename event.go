package cellnet

// 需要读取数据(tcp)
type ReadEvent struct {
	Ses Session
}

func (self *ReadEvent) Session() Session {
	return self.Ses
}

// 接收到数据(udp)
type RecvDataEvent struct {
	Ses  Session
	Data []byte
}

func (self *RecvDataEvent) Session() Session {
	return self.Ses
}

// 接收到消息
type RecvMsgEvent struct {
	Ses Session
	Msg interface{}
}

func (self *RecvMsgEvent) Session() Session {
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
	Ses Session
	Msg interface{} // 用户需要发送的消息
}

func (self *SendMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *SendMsgEvent) Session() Session {
	return self.Ses
}
