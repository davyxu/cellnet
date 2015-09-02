package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

// Connector/Acceptor事件
type SocketConnectError struct {
	Err error
}

type SocketListenError struct {
	Err error
}

const (
	SessionAccepted  = 1
	SessionConnected = 2
)

type SessionCreateType int

type SocketCreateSession struct {
	Stream cellnet.PacketStream
	Type   SessionCreateType
}

// Session相关事件

type SocketNewSession struct {
	Session cellnet.CellID
	Type    SessionCreateType
}

func (self SocketNewSession) GetSession() cellnet.CellID {
	return self.Session
}

type SocketData struct {
	Session cellnet.CellID
	Packet  *cellnet.Packet
}

func (self SocketData) GetSession() cellnet.CellID {
	return self.Session
}

func (self SocketData) GetPacket() *cellnet.Packet {
	return self.Packet
}

type SocketClose struct {
	Session cellnet.CellID
	Err     error
}

func (self SocketClose) GetSession() cellnet.CellID {
	return self.Session
}
