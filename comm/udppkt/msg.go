package udppkt

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
	cellnet.RegisterMessageMeta("binary", "udppkt.SessionAccepted", reflect.TypeOf((*SessionAccepted)(nil)).Elem(), 64001)
	cellnet.RegisterMessageMeta("binary", "udppkt.SessionConnected", reflect.TypeOf((*SessionConnected)(nil)).Elem(), 64002)
	cellnet.RegisterMessageMeta("binary", "udppkt.SessionConnectError", reflect.TypeOf((*SessionConnectError)(nil)).Elem(), 64003)
	cellnet.RegisterMessageMeta("binary", "udppkt.SessionClosed", reflect.TypeOf((*SessionClosed)(nil)).Elem(), 64004)
}
