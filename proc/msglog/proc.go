package msglog

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/comm"
)

type nameAddressGetter interface {
	Name() string
	Address() string
}

type LogHooker struct {
}

func (LogHooker) OnInboundEvent(ev cellnet.Event) {

	msg := ev.Message()
	ses := ev.Session()

	if IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(nameAddressGetter)

	switch msg := msg.(type) {
	case *cellnet.SessionAccepted:
		log.Debugf("#accepted(%s)@%d", peerInfo.Name(), ses.ID())
	case *cellnet.SessionClosed:
		log.Debugf("#closed(%s)@%d | Reason: %s", peerInfo.Name(), ses.ID(), msg.Error)
	case *cellnet.SessionConnected:
		log.Debugf("#connected(%s)@%d", peerInfo.Name(), ses.ID())
	case *cellnet.SessionConnectError:
		log.Debugf("#connectfailed(%s)@%d address: %s", peerInfo.Name(), ses.ID(), peerInfo.Address())
	default:
		log.Debugf("#recv(%s)@%d %s(%d) | %s",
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageToName(msg),
			cellnet.MessageToID(msg),
			cellnet.MessageToString(msg))
	}
}
func (LogHooker) OnOutboundEvent(ev cellnet.Event) {

	msg := ev.Message()
	ses := ev.Session()

	if rawPkt, ok := msg.(comm.RawPacket); ok {
		rawMsg, _, err := codec.DecodeMessage(rawPkt.MsgID, rawPkt.MsgData)
		if err != nil {
			log.Errorf("process msg log decode error: %s", err)
			return
		}

		msg = rawMsg
	}

	if IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(nameAddressGetter)

	log.Debugf("#send(%s)@%d %s(%d) | %s",
		peerInfo.Name(),
		ses.ID(),
		cellnet.MessageToName(msg),
		cellnet.MessageToID(msg),
		cellnet.MessageToString(msg))
}
