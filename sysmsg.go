package cellnet

import "fmt"

type SessionInit struct {
}

type SessionAccepted struct {
}

type SessionConnected struct {
}

type SessionConnectError struct {
}

type CloseReason int32

const (
	CloseReason_IO     CloseReason = iota // 普通IO断开
	CloseReason_Manual                    // 关闭前，调用过Session.Close
)

func (self CloseReason) String() string {
	switch self {
	case CloseReason_IO:
		return "IO"
	case CloseReason_Manual:
		return "Manual"
	}

	return "Unknown"
}

type SessionClosed struct {
	Reason CloseReason // 断开原因
}

// udp通知关闭,内部使用
type SessionCloseNotify struct {
}

func (self *SessionInit) String() string         { return fmt.Sprintf("%+v", *self) }
func (self *SessionAccepted) String() string     { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnected) String() string    { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnectError) String() string { return fmt.Sprintf("%+v", *self) }
func (self *SessionClosed) String() string       { return fmt.Sprintf("%+v", *self) }
func (self *SessionCloseNotify) String() string  { return fmt.Sprintf("%+v", *self) }

// 标记系统消息
func (self *SessionInit) SystemMessage()         {}
func (self *SessionAccepted) SystemMessage()     {}
func (self *SessionConnected) SystemMessage()    {}
func (self *SessionConnectError) SystemMessage() {}
func (self *SessionClosed) SystemMessage()       {}
func (self *SessionCloseNotify) SystemMessage()  {}

// 使用类型断言判断是否为系统消息
type SystemMessageIdentifier interface {
	SystemMessage()
}
