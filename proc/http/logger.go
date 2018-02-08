package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc/msglog"
)

type LogHooker struct {
}

func (LogHooker) OnInboundEvent(raw cellnet.Event) {

	msg := raw.Message()
	ses := raw.BaseSession()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch raw.(type) {
	case *cellnet.RecvMsgEvent:
		log.Debugf("#http.%s(%s) %s %s | %s",
			httpContext.Request().Method,
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}

}
func (LogHooker) OnOutboundEvent(raw cellnet.Event) {

	msg := raw.Message()
	ses := raw.BaseSession()

	if msglog.IsBlockedMessageByID(cellnet.MessageToID(msg)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	httpContext := ses.(HttpContext)

	switch raw.(type) {
	case *cellnet.SendMsgEvent:
		log.Debugf("#http.Respond(%s) %s %s | %s",
			peerInfo.Name(),
			httpContext.Request().URL.Path,
			cellnet.MessageToName(msg),
			cellnet.MessageToString(msg))
	}
}
