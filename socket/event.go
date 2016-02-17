package socket

import (
	"github.com/davyxu/cellnet"
)

var (
	Event_SessionConnected = uint32(cellnet.MessageMetaByName("coredef.SessionConnected").ID)
	Event_SessionClosed    = uint32(cellnet.MessageMetaByName("coredef.SessionClosed").ID)
	Event_SessionAccepted  = uint32(cellnet.MessageMetaByName("coredef.SessionAccepted").ID)
	Event_PeerInit         = uint32(cellnet.MessageMetaByName("coredef.PeerInit").ID)
	Event_PeerStart        = uint32(cellnet.MessageMetaByName("coredef.PeerStart").ID)
	Event_PeerStop         = uint32(cellnet.MessageMetaByName("coredef.PeerStop").ID)
)

// 会话事件
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

// 端事件
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
