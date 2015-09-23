package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
)

var (
	Event_Connected = uint32(cellnet.Type2ID(&coredef.ConnectedACK{}))
	Event_Closed    = uint32(cellnet.Type2ID(&coredef.ClosedACK{}))
	Event_Accepted  = uint32(cellnet.Type2ID(&coredef.AcceptedACK{}))
)

type DataEvent struct {
	*cellnet.Packet
	Ses cellnet.Session
}

func NewDataEvent(msgid uint32, s cellnet.Session, data []byte) *DataEvent {
	return &DataEvent{
		Packet: &cellnet.Packet{MsgID: msgid, Data: data},
		Ses:    s,
	}
}

type peerProfile struct {
	queue *cellnet.EvQueue
	name  string
}

func (self *peerProfile) SetName(name string) {
	self.name = name
}

func (self *peerProfile) Name() string {
	return self.name
}
