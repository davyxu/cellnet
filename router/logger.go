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
	MessageId() int
	Message() any
}

func RecvLogger(where, data any) {

	var msg any

	switch v := data.(type) {
	case MessageFetcher:
		if cellmsglog.IsMsgVisible(v.MessageId()) {
			msg = v.Message()
		} else {
			return
		}
	default:
		if cellmsglog.IsMsgVisible(cellmeta.MessageToId(v)) {
			msg = v
		} else {
			return
		}
	}

	xlog.Debugf("#recv %v len: %d | %s %s", where, cellmeta.MessageSize(msg), cellmeta.MessageToName(msg), cellmeta.MessageToString(msg))
}

func SendLogger(where any, msg any) {

	switch v := msg.(type) {
	case cellevent.SystemMessageIdentifier:
		return
	case cellmsglog.PacketMessagePeeker:
		msg = v.Message()
	case *cellpeer.RawPacket:
		msg, _, _ = cellcodec.Decode(v.MsgId, v.MsgData)
	}

	if cellmsglog.IsMsgVisible(cellmeta.MessageToId(msg)) {
		xlog.Debugf("#send %v len: %d | %s %s", where, cellmeta.MessageSize(msg), cellmeta.MessageToName(msg), cellmeta.MessageToString(msg))
	}
}
