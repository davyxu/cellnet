package msglog

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
)

// 萃取消息中的消息
type PacketMessagePeeker interface {
	Message() interface{}
}

func WriteRecvLogger(log *golog.Logger, protocol string, ses cellnet.Session, msg interface{}) {

	if log.IsDebugEnabled() {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnet.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnet.PeerProperty)

			log.Debugf("#%s.recv(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				cellnet.MessageSize(msg),
				cellnet.MessageToName(msg),
				cellnet.MessageToString(msg))
		}

	}
}

func WriteSendLogger(log *golog.Logger, protocol string, ses cellnet.Session, msg interface{}) {

	if log.IsDebugEnabled() {

		if peeker, ok := msg.(PacketMessagePeeker); ok {
			msg = peeker.Message()
		}

		if IsMsgLogValid(cellnet.MessageToID(msg)) {
			peerInfo := ses.Peer().(cellnet.PeerProperty)

			log.Debugf("#%s.send(%s)@%d len: %d %s | %s",
				protocol,
				peerInfo.Name(),
				ses.ID(),
				cellnet.MessageSize(msg),
				cellnet.MessageToName(msg),
				cellnet.MessageToString(msg))
		}

	}

}
