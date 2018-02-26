package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc/msglog"
)

type LogHooker struct {
}

func (LogHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msg := inputEvent.Message()
	ses := inputEvent.Session()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch inputEvent.(type) {
	case *cellnet.RecvMsgEvent:
		log.Debugf("#recv %s(%s) %s %s | %s",
			httpContext.Request().Method,
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

	return inputEvent

}
func (LogHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msg := inputEvent.Message()
	ses := inputEvent.Session()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch inputEvent.(type) {
	case *cellnet.SendMsgEvent:
		log.Debugf("#send Response(%s) %s %s | %s",
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

	return inputEvent
}
