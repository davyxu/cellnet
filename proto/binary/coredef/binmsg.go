package coredef

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
)

type SessionAccepted struct {
}

type SessionConnected struct {
}

type SessionAcceptFailed struct {
	Result cellnet.Result
}

type SessionConnectFailed struct {
	Result cellnet.Result
}

type SessionClosed struct {
	Result cellnet.Result
}

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}
