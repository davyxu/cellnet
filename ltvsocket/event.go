package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

const (
	Event_RecvData  = 1
	Event_Connected = 2
	Event_Closed    = 3
	Event_Accepted  = 4
)

type DataEvent struct {
	*cellnet.Packet
	Ses Session
}

func NewDataEvent(msgid uint32, s Session, data []byte) *DataEvent {
	return &DataEvent{
		Packet: &cellnet.Packet{MsgID: msgid, Data: data},
		Ses:    s,
	}
}
