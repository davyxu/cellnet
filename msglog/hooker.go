package cellmsglog

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	cellpeer "github.com/davyxu/cellnet/peer"
	"github.com/davyxu/xlog"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

var (
	EnableMsgLog     = true
	SystemMsgVisible = true
)

func RecvLogger(input *cellevent.RecvMsg) *cellevent.RecvMsg {

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

			// blue
			xlog.Debugf("#recv %d %s %d %s",
				getSessionID(input.Ses),
				cellmeta.MessageToName(msg),
				cellmeta.MessageSize(msg),
				cellmeta.MessageToString(msg))
		}

	}

	return input
}

func getSessionID(session cellnet.Session) int64 {

	if fetcher, ok := session.(cellpeer.SessionID64Fetcher); ok {
		return fetcher.ID()
	}
	return 0
}

func SendLogger(input *cellevent.SendMsg) *cellevent.SendMsg {

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

			// purple
			xlog.Debugf("#send %d %s %d %s",
				getSessionID(input.Ses),
				cellmeta.MessageToName(msg),
				cellmeta.MessageSize(msg),
				cellmeta.MessageToString(msg))
		}

	}

	return input
}
