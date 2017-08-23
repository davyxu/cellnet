package coredef

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
)

// Acceptor的会话被允许接入时
type SessionAccepted struct {
}

// Connector的连接创建时
type SessionConnected struct {
}

// Acceptor的会话接受错误
type SessionAcceptFailed struct {
	Result cellnet.Result
}

// Connector连接失败时
type SessionConnectFailed struct {
	Result cellnet.Result
}

// Session连接断开时
type SessionClosed struct {
	Result cellnet.Result
}

// 内部消息,勿使用及注册响应
type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}
