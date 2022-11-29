package cellevent

import (
	"fmt"
	"github.com/davyxu/cellnet"
	cellmeta "github.com/davyxu/cellnet/meta"
	"reflect"
)

type SessionAccepted struct {
}

type SessionConnected struct {
}

type SessionConnectError struct {
	Err            error
	ConnectedTimes int32
	RetryTimes     int32
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
	Err    error       // 非EOF和网络读取错误
}

// udp通知关闭,内部使用
type SessionCloseNotify struct {
}

func (self *SessionAccepted) String() string     { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnected) String() string    { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnectError) String() string { return fmt.Sprintf("%+v", *self) }
func (self *SessionClosed) String() string       { return fmt.Sprintf("%+v", *self) }
func (self *SessionCloseNotify) String() string  { return fmt.Sprintf("%+v", *self) }

// 标记系统消息
func (self *SessionAccepted) SystemMessage()     {}
func (self *SessionConnected) SystemMessage()    {}
func (self *SessionConnectError) SystemMessage() {}
func (self *SessionClosed) SystemMessage()       {}
func (self *SessionCloseNotify) SystemMessage()  {}

// 使用类型断言判断是否为系统消息
type SystemMessageIdentifier interface {
	SystemMessage()
}

func BuildSystemEvent(ses cellnet.Session, msg any) *RecvMsg {

	meta := cellmeta.MetaByMsg(msg)
	if meta == nil {
		panic("sysmsg meta not found")
	}

	return &RecvMsg{
		Ses:   ses,
		MsgId: meta.Id,
		Msg:   msg,
	}
}

func init() {
	cellmeta.Register(&cellmeta.Meta{
		Type: reflect.TypeOf((*SessionAccepted)(nil)).Elem(),
		Id:   1,
	})
	cellmeta.Register(&cellmeta.Meta{
		Type: reflect.TypeOf((*SessionConnected)(nil)).Elem(),
		Id:   2,
	})
	cellmeta.Register(&cellmeta.Meta{
		Type: reflect.TypeOf((*SessionConnectError)(nil)).Elem(),
		Id:   3,
	})
	cellmeta.Register(&cellmeta.Meta{
		Type: reflect.TypeOf((*SessionClosed)(nil)).Elem(),
		Id:   4,
	})
	cellmeta.Register(&cellmeta.Meta{
		Type: reflect.TypeOf((*SessionCloseNotify)(nil)).Elem(),
		Id:   5,
	})
}
