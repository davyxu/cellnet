package cellnet

import (
	"fmt"
	"reflect"
)

type EventHandler interface {
	Call(*SessionEvent)
}

var EnableHandlerLog bool

// 显示handler的名称
func HandlerName(h EventHandler) string {
	if h == nil {
		return "nil"
	}

	return reflect.TypeOf(h).Elem().Name()
}

type HandlerOptional struct {
	Handler EventHandler
	Enable  bool
}

// 链接一连串handler, 返回第一个
func HandlerLink(rawList ...interface{}) (ret []EventHandler) {

	for _, raw := range rawList {
		switch v := raw.(type) {
		case EventHandler:
			ret = append(ret, v)
		case HandlerOptional:
			v = raw.(HandlerOptional)
			if v.Enable {
				ret = append(ret, v.Handler)
			}
		case []EventHandler:
			ret = append(ret, v...)
		default:
			panic("Require 'EventHandler', 'HandlerOptional', []EventHandler: " + fmt.Sprintln(reflect.TypeOf(raw)))
		}
	}

	return
}

func HandlerChainListName(hlist []EventHandler) {

	for _, h := range hlist {

		if EnableHandlerLog {
			log.Debugf("%s", HandlerName(h))
		}

	}

}

func HandlerChainCall(hlist []EventHandler, ev *SessionEvent) {

	for _, h := range hlist {

		if EnableHandlerLog {
			log.Debugf("%d %s [%s] <%s> MsgID: %d(%s) {%s} Raw: (%d)%v Tag: %v TransmitTag: %v", ev.UID, ev.TypeString(), ev.PeerName(), HandlerName(h), ev.MsgID, ev.MsgName(), ev.MsgString(), ev.MsgSize(), ev.Data, ev.Tag, ev.TransmitTag)
		}

		h.Call(ev)
	}

}
