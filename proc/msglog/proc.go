package msglog

import (
	"github.com/davyxu/cellnet"
)

func WriteRecvLogger(ses cellnet.Session, msg interface{}) {

	if log.IsDebugEnabled() && !IsBlockedMessageByID(cellnet.MessageToID(msg)) {

		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.Debugf("#recv(%s)@%d %s(%d) | %s",
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageToName(msg),
			cellnet.MessageToID(msg),
			cellnet.MessageToString(msg))
	}
}

func WriteSendLogger(ses cellnet.Session, msg interface{}) {

	if log.IsDebugEnabled() && !IsBlockedMessageByID(cellnet.MessageToID(msg)) {

		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.Debugf("#send(%s)@%d %s(%d) | %s",
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageToName(msg),
			cellnet.MessageToID(msg),
			cellnet.MessageToString(msg))
	}

}
