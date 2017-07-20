package cellnet

import (
	"fmt"
	"reflect"
)

type EventHandler interface {
	Call(*Event)
}

var EnableHandlerLog bool

// 显示handler的名称
func HandlerName(h EventHandler) string {
	if h == nil {
		return "nil"
	}

	return reflect.TypeOf(h).Elem().Name()
}

func HandlerString(h EventHandler) string {

	if sg, ok := h.(fmt.Stringer); ok {
		return sg.String()
	} else {
		return HandlerName(h)
	}
}

func HandlerLog(h EventHandler, ev *Event) {

	if EnableHandlerLog {
		log.Debugf("evid: %d #%s [%s] chain: %d <%s> SesID: %d Result: %d MsgID: %d(%s) {%s} Tag: %v TransmitTag: %v Raw: (%d)%v",
			ev.UID,
			ev.Type.String(),
			ev.PeerName(),
			ev.chainid,
			HandlerString(h),
			ev.SessionID(),
			ev.Result(),
			ev.MsgID,
			ev.MsgName(),
			ev.MsgString(),
			ev.Tag,
			ev.TransmitTag,
			ev.MsgSize(),
			ev.Data,
		)
	}
}

func HandlerChainCall(hlist []EventHandler, ev *Event) {

	for _, h := range hlist {

		HandlerLog(h, ev)

		h.Call(ev)

		if ev.Result() != Result_OK {
			break
		}
	}

}
