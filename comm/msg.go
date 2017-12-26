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
}

func (self *SessionAccepted) String() string     { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnected) String() string    { return fmt.Sprintf("%+v", *self) }
func (self *SessionConnectError) String() string { return fmt.Sprintf("%+v", *self) }
func (self *SessionClosed) String() string       { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta("binary", "tcpproc.SessionAccepted", reflect.TypeOf((*SessionAccepted)(nil)).Elem(), 63001)
	cellnet.RegisterMessageMeta("binary", "tcpproc.SessionConnected", reflect.TypeOf((*SessionConnected)(nil)).Elem(), 63002)
	cellnet.RegisterMessageMeta("binary", "tcpproc.SessionConnectError", reflect.TypeOf((*SessionConnectError)(nil)).Elem(), 63003)
	cellnet.RegisterMessageMeta("binary", "tcpproc.SessionClosed", reflect.TypeOf((*SessionClosed)(nil)).Elem(), 63004)
}
