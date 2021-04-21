package cellmsglog

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	"github.com/davyxu/ulog"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

var (
	EnableMsgLog     = true
	SystemMsgVisible = true
)

func RecvLogger(input *cellevent.RecvMsgEvent) *cellevent.RecvMsgEvent {

	if EnableMsgLog {

		msg := input.Message()
		msgID := input.MessageID()

		if !SystemMsgVisible {
			if _, ok := msg.(cellevent.SystemMessageIdentifier); ok {
				return input
			}
		}

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgVisible(msgID) {

			sesID := input.Session().(interface {
				ID() int64
			}).ID()

			ulog.Debugf("#recv %d %s %d %s",
				sesID,
				cellmeta.MessageToName(msg),
				cellmeta.MessageSize(msg),
				cellmeta.MessageToString(msg))
		}

	}

	return input
}

func SendLogger(input *cellevent.SendMsgEvent) *cellevent.SendMsgEvent {

	if EnableMsgLog {

		msg := input.Message()
		msgID := input.MessageID()

		if !SystemMsgVisible {
			if _, ok := msg.(cellevent.SystemMessageIdentifier); ok {
				return input
			}
		}

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgVisible(msgID) {

			sesID := input.Session().(interface {
				ID() int64
			}).ID()

			ulog.Debugf("#send %d %s %d %s",
				sesID,
				cellmeta.MessageToName(msg),
				cellmeta.MessageSize(msg),
				cellmeta.MessageToString(msg))
		}

	}

	return input
}
