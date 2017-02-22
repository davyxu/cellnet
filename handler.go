package cellnet

import "reflect"

type EventHandler interface {
	Call(*SessionEvent)

	Next() EventHandler
	SetNext(EventHandler) EventHandler

	SetTag(interface{})
	Tag() interface{}
	MatchTag(interface{}) bool
}

type BaseEventHandler struct {
	next EventHandler

	tag interface{}
}

func (self *BaseEventHandler) SetTag(t interface{}) {
	self.tag = t
}

func (self *BaseEventHandler) Tag() interface{} {
	return self.tag
}

func (self *BaseEventHandler) MatchTag(t interface{}) bool {
	return self.tag == t
}

func (self *BaseEventHandler) Next() EventHandler {
	return self.next
}

func (self *BaseEventHandler) SetNext(next EventHandler) EventHandler {
	self.next = next
	return next
}

var EnableHandlerLog bool

// 显示handler的名称
func HandlerName(h EventHandler) string {
	if h == nil {
		return "nil"
	}

	return reflect.TypeOf(h).Elem().Name()
}

// 链接一连串handler, 返回第一个
func LinkHandler(hlist ...EventHandler) EventHandler {

	var pre EventHandler

	for _, h := range hlist {

		if h == nil {
			continue
		}

		if pre != nil {
			pre.SetNext(h)
		}

		pre = h
	}

	if len(hlist) == 0 {
		return nil
	}

	return hlist[0]
}

func HandlerChainListName(h EventHandler) {

	for h != nil {

		if EnableHandlerLog {
			log.Debugf("%s", HandlerName(h))
		}

		h = h.Next()
	}

}

func HandlerChainCall(h EventHandler, ev *SessionEvent) {

	for h != nil {

		if EnableHandlerLog {
			log.Debugf("%d %s [%s] <%s> MsgID: %d(%s) {%s} Raw: (%d)%v Tag: %v TransmitTag: %v", ev.UID, ev.TypeString(), ev.PeerName(), HandlerName(h), ev.MsgID, ev.MsgName(), ev.MsgString(), ev.MsgSize(), ev.Data, ev.Tag, ev.TransmitTag)
		}

		h.Call(ev)

		h = h.Next()
	}

}
