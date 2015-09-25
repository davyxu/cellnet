package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
)

var (
	Event_SessionConnected = uint32(cellnet.Type2ID(&coredef.SessionConnected{}))
	Event_SessionClosed    = uint32(cellnet.Type2ID(&coredef.SessionClosed{}))
	Event_SessionAccepted  = uint32(cellnet.Type2ID(&coredef.SessionAccepted{}))
	Event_PeerInit         = uint32(cellnet.Type2ID(&coredef.PeerInit{}))
	Event_PeerStart        = uint32(cellnet.Type2ID(&coredef.PeerStart{}))
	Event_PeerStop         = uint32(cellnet.Type2ID(&coredef.PeerStop{}))
)

// 内部事件
type SessionEvent struct {
	*cellnet.Packet
	Ses cellnet.Session
}

func NewSessionEvent(msgid uint32, s cellnet.Session, data []byte) *SessionEvent {
	return &SessionEvent{
		Packet: &cellnet.Packet{MsgID: msgid, Data: data},
		Ses:    s,
	}
}

type PeerEvent struct {
	MsgID uint32
	P     cellnet.Peer
}

func (self PeerEvent) ContextID() int {
	return int(self.MsgID)
}

func NewPeerEvent(msgid uint32, p cellnet.Peer) *PeerEvent {
	return &PeerEvent{MsgID: msgid, P: p}
}
