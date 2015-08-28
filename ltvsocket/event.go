package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

type EventConnectError struct {
	Err error
}

type EventListenError struct {
	Err error
}

const (
	SessionAccepted  = 1
	SessionConnected = 2
)

type SessionCreateType int

type EventCreateSession struct {
	Stream cellnet.PacketStream
	Type   SessionCreateType
}

type EventNewSession struct {
	Session cellnet.CellID
	Type    SessionCreateType
}

type EventData struct {
	Session cellnet.CellID
	Packet  *cellnet.Packet
}

type EventClose struct {
	Session cellnet.CellID
	Err     error
}
