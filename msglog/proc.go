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
	EnableMsgLog = true
)

func WriteRecvLogger(protocol string, ses cellnet.Session, msg interface{}) {

	if EnableMsgLog {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnet.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnet.PeerProperty)

			ulog.Debugf("#%s.recv(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				cellnet.MessageSize(msg),
				cellnet.MessageToName(msg),
				cellnet.MessageToString(msg))
		}

	}
}

func WriteSendLogger(protocol string, ses cellnet.Session, msg interface{}) {

	if EnableMsgLog {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnet.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnet.PeerProperty)

			ulog.Debugf("#%s.send(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				cellnet.MessageSize(msg),
				cellnet.MessageToName(msg),
				cellnet.MessageToString(msg))
		}

	}

}

type Hooker struct {
}

// 来自后台服务器的消息
func (Hooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	WriteRecvLogger("tcp", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

// 发送给后台服务器
func (Hooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	WriteSendLogger("tcp", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}
