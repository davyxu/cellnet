package msglog

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

func WriteRecvLogger(protocol string, ses cellnet.Session, msg interface{}) {

	if ulog.IsLevelEnabled(ulog.DebugLevel) {

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

	if ulog.IsLevelEnabled(ulog.DebugLevel) {

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
