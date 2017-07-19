package coredef

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/davyxu/goobjfmt"
	"reflect"
)

type SessionAccepted struct {
}

func (m *SessionAccepted) String() string { return goobjfmt.CompactTextString(m) }

type SessionConnected struct {
}

func (m *SessionConnected) String() string { return goobjfmt.CompactTextString(m) }

type SessionAcceptFailed struct {
	Result cellnet.Result
}

func (m *SessionAcceptFailed) String() string { return goobjfmt.CompactTextString(m) }

type SessionConnectFailed struct {
	Result cellnet.Result
}

func (m *SessionConnectFailed) String() string { return goobjfmt.CompactTextString(m) }

type SessionClosed struct {
	Result cellnet.Result
}

func (m *SessionClosed) String() string { return goobjfmt.CompactTextString(m) }

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (m *RemoteCallACK) String() string { return goobjfmt.CompactTextString(m) }

func init() {

	// coredef.proto
	cellnet.RegisterMessageMeta("binary", "coredef.SessionAccepted", reflect.TypeOf((*SessionAccepted)(nil)).Elem(), 3495179174)
	cellnet.RegisterMessageMeta("binary", "coredef.SessionConnected", reflect.TypeOf((*SessionConnected)(nil)).Elem(), 3551021301)
	cellnet.RegisterMessageMeta("binary", "coredef.SessionAcceptFailed", reflect.TypeOf((*SessionAcceptFailed)(nil)).Elem(), 3277953230)
	cellnet.RegisterMessageMeta("binary", "coredef.SessionConnectFailed", reflect.TypeOf((*SessionConnectFailed)(nil)).Elem(), 3980285497)
	cellnet.RegisterMessageMeta("binary", "coredef.SessionClosed", reflect.TypeOf((*SessionClosed)(nil)).Elem(), 3480086952)
	cellnet.RegisterMessageMeta("binary", "coredef.RemoteCallACK", reflect.TypeOf((*RemoteCallACK)(nil)).Elem(), 2811469770)

}
