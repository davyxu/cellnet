package comm

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
	"reflect"
)

type SessionAccepted struct {
}

type SessionConnected struct {
}

type SessionConnectError struct {
}

type SessionClosed struct {
	Error string
}

// udp通知关闭
type SessionCloseNotify struct {
}

type RemoteCallREQ struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self *SessionAccepted) String() string     { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnected) String() string    { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnectError) String() string { return fmt.Sprintf("%+v", *self) }
func (self *SessionClosed) String() string       { return fmt.Sprintf("%+v", *self) }
func (self *SessionCloseNotify) String() string  { return fmt.Sprintf("%+v", *self) }
func (self *RemoteCallREQ) String() string       { return fmt.Sprintf("%+v", *self) }
func (self *RemoteCallACK) String() string       { return fmt.Sprintf("%+v", *self) }

func (self *RemoteCallREQ) GetMsgID() uint32   { return self.MsgID }
func (self *RemoteCallREQ) GetMsgData() []byte { return self.Data }
func (self *RemoteCallREQ) GetCallID() int64   { return self.CallID }
func (self *RemoteCallACK) GetMsgID() uint32   { return self.MsgID }
func (self *RemoteCallACK) GetMsgData() []byte { return self.Data }
func (self *RemoteCallACK) GetCallID() int64   { return self.CallID }

func init() {
	cellnet.RegisterMessageMeta("binary", "comm.SessionAccepted", reflect.TypeOf((*SessionAccepted)(nil)).Elem(), 63001)
	cellnet.RegisterMessageMeta("binary", "comm.SessionConnected", reflect.TypeOf((*SessionConnected)(nil)).Elem(), 63002)
	cellnet.RegisterMessageMeta("binary", "comm.SessionConnectError", reflect.TypeOf((*SessionConnectError)(nil)).Elem(), 63003)
	cellnet.RegisterMessageMeta("binary", "comm.SessionClosed", reflect.TypeOf((*SessionClosed)(nil)).Elem(), 63004)
	cellnet.RegisterMessageMeta("binary", "comm.SessionCloseNotify", reflect.TypeOf((*SessionCloseNotify)(nil)).Elem(), 63005)
	cellnet.RegisterMessageMeta("binary", "comm.RemoteCallREQ", reflect.TypeOf((*RemoteCallREQ)(nil)).Elem(), 63006)
	cellnet.RegisterMessageMeta("binary", "comm.RemoteCallACK", reflect.TypeOf((*RemoteCallACK)(nil)).Elem(), 63007)
}
