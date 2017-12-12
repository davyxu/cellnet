package cellnet

// 需要读取数据
type ReadEvent struct {
	Ses Session
}

// 接收到数据
type RecvDataEvent struct {
	Ses  Session
	Data []byte
}

// 接收到消息
type RecvMsgEvent struct {
	Ses Session
	Msg interface{}
}

func (self *RecvMsgEvent) Session() Session {
	return self.Ses
}

func (self *RecvMsgEvent) GetMsg() interface{} {
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

func (self *SendMsgEvent) GetMsg() interface{} {
	return self.Msg
}

func (self *SendMsgEvent) Session() Session {
	return self.Ses
}

// 会话接收数据时发生错误的事件
type RecvErrorEvent struct {
	Ses   Session
	Error error
}

// 会话发送数据时发生错误的事件
type SendMsgErrorEvent struct {
	Ses   Session
	Error error
	Msg   interface{}
}

// 连接错误事件
type SessionConnectErrorEvent struct {
	Ses   Session
	Error error
}

// 会话连接关闭事件
type SessionClosedEvent struct {
	Ses   Session
	Error error
}

// 已连接上远方服务器事件
type SessionConnectedEvent struct {
	Ses Session
}

// 已接受一个连接事件
type SessionAcceptedEvent struct {
	Ses Session
}
