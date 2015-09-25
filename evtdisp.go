package cellnet

import (
	"container/list"
	"reflect"
)

type EventDispatcher struct {
	handler map[string]*list.List
}

func (self *EventDispatcher) Add(name string, callback func(...interface{})) {
	arr := self.handler[name]

	if arr == nil {
		arr = list.New()

		self.handler[name] = arr
	}

	arr.PushBack(callback)
}

func (self *EventDispatcher) Invoke(name string, args ...interface{}) {
	if v, ok := self.handler[name]; ok {

		for e := v.Front(); e != nil; e = e.Next() {
			c := e.Value.(func(...interface{}))
			c(args...)
		}
	}

}

func (self *EventDispatcher) Remove(name string, callback func(...interface{})) {
	arr := self.handler[name]

	for e := arr.Front(); e != nil; e = e.Next() {

		c := e.Value.(func(...interface{}))

		if reflect.ValueOf(c).Pointer() == reflect.ValueOf(callback).Pointer() {
			arr.Remove(e)
		}

	}

}

func (self *EventDispatcher) Clear() {
	self.handler = make(map[string]*list.List)
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handler: make(map[string]*list.List),
	}
}
