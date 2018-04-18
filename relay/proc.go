package relay

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

// 处理入站的relay消息
func ResoleveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		userMsg, _, err := codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Data)
		if err == nil {

			if log.IsDebugEnabled() {

				peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

				log.Debugf("#relay.recv(%s)@%d len: %d %s context: %v | %s",
					peerInfo.Name(),
					inputEvent.Session().ID(),
					cellnet.MessageSize(userMsg),
					cellnet.MessageToName(userMsg),
					relayMsg.ContextID,
					cellnet.MessageToString(userMsg))
			} else {
				log.Errorln("relay.ResoleveInboundEvent:", err)
			}

			ev := &RecvMsgEvent{
				inputEvent.Session(),
				userMsg,
				relayMsg.ContextID,
			}

			// 转到对应线程中调用
			cellnet.SessionQueuedCall(inputEvent.Session(), func() {

				broadcast(ev)
			})

			ouputEvent = ev
			handled = true

			return
		}
	}

	return inputEvent, false
}

// 处理relay.Relay出站消息的日志
func ResolveOutboundEvent(inputEvent cellnet.Event) (handled bool) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		if log.IsDebugEnabled() {

			userMsg, _, err := codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Data)
			if err == nil {

				peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

				log.Debugf("#relay.send(%s)@%d len: %d %s context: %v | %s",
					peerInfo.Name(),
					inputEvent.Session().ID(),
					cellnet.MessageSize(userMsg),
					cellnet.MessageToName(userMsg),
					relayMsg.ContextID,
					cellnet.MessageToString(userMsg))

				return true
			}

		}

	}

	return
}
