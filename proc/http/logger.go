package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc/msglog"
)

type LogHooker struct {
}

func (LogHooker) OnInboundEvent(raw cellnet.Event) {

	msg := raw.Message()
	ses := raw.Session()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch raw.(type) {
	case *cellnet.RecvMsgEvent:
		log.Debugf("#recv %s(%s) %s %s | %s",
			httpContext.Request().Method,
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}
func (LogHooker) OnOutboundEvent(raw cellnet.Event) {

	msg := raw.Message()
	ses := raw.Session()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch raw.(type) {
	case *cellnet.SendMsgEvent:
		log.Debugf("#send Response(%s) %s %s | %s",
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}
}
