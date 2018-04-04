package relay

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

func ResoleveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		userMsg, meta, err := codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Data)
		if err == nil {

			if log.IsDebugEnabled() {

				peerInfo := inputEvent.Session().Peer().(interface {
					Name() string
				})

				log.Debugf("#relay.recv(%s)@%d %s(%d) | %s",
					peerInfo.Name(),
					inputEvent.Session().ID(),
					meta.TypeName(),
					meta.ID,
					cellnet.MessageToString(userMsg))
			}

			ev := &RecvMsgEvent{
				inputEvent.Session(),
				userMsg,
				relayMsg.SessionID,
			}

			if bcFunc == nil || bcFunc(ev) {
				ouputEvent = ev
				handled = true
			}

			return
		}
	}

	return inputEvent, false
}

func ResolveOutboundEvent(inputEvent cellnet.Event) (handled bool) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		userMsg, meta, err := codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Data)
		if err == nil {

			peerInfo := inputEvent.Session().Peer().(interface {
				Name() string
			})

			log.Debugf("#relay.send(%s)@%d %s(%d) | %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				meta.TypeName(),
				meta.ID,
				cellnet.MessageToString(userMsg))

			return true
		}
	}

	return
}
