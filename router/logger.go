package cellrouter

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	cellmsglog "github.com/davyxu/cellnet/msglog"
	cellpeer "github.com/davyxu/cellnet/peer"
	"github.com/davyxu/xlog"
)

type MessageFetcher interface {
	MessageID() int
	Message() interface{}
}

func RecvLogger(where interface{}, data interface{}) {

	var msg interface{}

	switch v := data.(type) {
	case MessageFetcher:
		if cellmsglog.IsMsgVisible(v.MessageID()) {
			msg = v.Message()
		} else {
			return
		}
	default:
		if cellmsglog.IsMsgVisible(cellmeta.MessageToID(v)) {
			msg = v
		} else {
			return
		}
	}

	xlog.Debugf("#recv %v len: %d | %s %s", where, cellmeta.MessageSize(msg), cellmeta.MessageToName(msg), cellmeta.MessageToString(msg))
}

func SendLogger(where interface{}, msg interface{}) {

	switch v := msg.(type) {
	case cellevent.SystemMessageIdentifier:
		return
	case cellmsglog.PacketMessagePeeker:
		msg = v.Message()
	case *cellpeer.RawPacket:
		msg, _, _ = cellcodec.Decode(v.MsgID, v.MsgData)
	}

	if cellmsglog.IsMsgVisible(cellmeta.MessageToID(msg)) {
		xlog.Debugf("#send %v len: %d | %s %s", where, cellmeta.MessageSize(msg), cellmeta.MessageToName(msg), cellmeta.MessageToString(msg))
	}
}
