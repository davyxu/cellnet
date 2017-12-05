package cellnet

// 会话开始接收数据事件
type RecvEvent struct {
	Ses Session
}

// 会话开始发送数据事件
type SendEvent struct {
	Ses Session
	Msg interface{} // 用户需要发送的消息
}

// 会话接收数据时发生错误的事件
type RecvErrorEvent struct {
	Ses   Session
	Error error
}

// 会话发送数据时发生错误的事件
type SendErrorEvent struct {
	Ses   Session
	Error error
	Msg   interface{}
}

// 连接错误事件
type ConnectErrorEvent struct {
	Ses   Session
	Error error
}

// 会话
type SessionCleanupEvent struct {
	Ses Session
}

// 会话连接关闭事件
type SessionClosedEvent struct {
	Ses   Session
	Error error
}

// 会话开始事件
type SessionStartEvent struct {
	Ses Session
}

// 已连接上远方服务器事件
type ConnectedEvent = SessionStartEvent

// 已接受一个连接事件
type AcceptedEvent = SessionStartEvent

// 会话退出事件
type SessionExitEvent struct {
	Ses Session
}
