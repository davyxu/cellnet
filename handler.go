package cellnet

import "reflect"

type EventHandler interface {
	Call(*SessionEvent) error

	SetNext(EventHandler) EventHandler
	Next() EventHandler

	SetTag(interface{})
	Tag() interface{}
	MatchTag(interface{}) bool
}

type BaseEventHandler struct {
	next EventHandler
	tag  interface{}
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

func (self *BaseEventHandler) SetNext(h EventHandler) EventHandler {
	self.next = h
	return h
}

func (self *BaseEventHandler) Next() EventHandler {
	return self.next
}

func (self *BaseEventHandler) CallNext(ev *SessionEvent) error {

	return HandlerCallNext(self.next, ev)
}

func HandlerCallFirst(h EventHandler, ev *SessionEvent) error {
	if EnableHandlerLog {
		log.Debugf("HandlerFirst: %s %s", HandlerName(h), ev.String())
	}

	return h.Call(ev)
}

func HandlerCallNext(h EventHandler, ev *SessionEvent) error {
	if EnableHandlerLog {
		log.Debugf("HandlerNext: %s %s", HandlerName(h), ev.String())
	}

	if h == nil {
		return nil
	}

	return h.Call(ev)
}

var EnableHandlerLog bool

// 显示handler的名称
func HandlerName(h EventHandler) string {
	if h == nil {
		return "nil"
	}

	return reflect.TypeOf(h).Elem().Name()
}

// handler的类型
func HandlerType(h EventHandler) reflect.Type {
	return reflect.TypeOf(h).Elem()
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
