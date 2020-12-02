package msglog

import (
	"github.com/davyxu/cellnet"
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

func WriteRecvLogger(protocol string, ev cellnet.Event) {

	if EnableMsgLog {

		msg := ev.Message()

		if !SystemMsgVisible {
			if _, ok := msg.(cellnet.SystemMessageIdentifier); ok {
				return
			}
		}

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnet.MessageToID(msg)) {
			peerInfo := ev.Session().Peer().(cellnet.PeerProperty)

			meta := cellnet.MessageMetaByID(ev.MessageID())

			ulog.Debugf("#%s.recv(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ev.Session().ID(),
				cellnet.MessageSize(msg),
				meta.TypeName(),
				cellnet.MessageToString(msg))
		}

	}
}

func WriteSendLogger(protocol string, ev cellnet.Event) {

	if EnableMsgLog {

		msg := ev.Message()

		if !SystemMsgVisible {
			if _, ok := msg.(cellnet.SystemMessageIdentifier); ok {
				return
			}
		}

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(ev.MessageID()) {
			peerInfo := ev.Session().Peer().(cellnet.PeerProperty)

			meta := cellnet.MessageMetaByID(ev.MessageID())

			ulog.Debugf("#%s.send(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ev.Session().ID(),
				cellnet.MessageSize(msg),
				meta.TypeName(),
				cellnet.MessageToString(msg))
		}

	}

}

type Hooker struct {
}

// 来自后台服务器的消息
func (Hooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	WriteRecvLogger("tcp", inputEvent)

	return inputEvent
}

// 发送给后台服务器
func (Hooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	WriteSendLogger("tcp", inputEvent)

	return inputEvent
}
