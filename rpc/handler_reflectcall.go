package rpc

import (
	"reflect"

	"github.com/davyxu/cellnet"
)

type ReflectCallHandler struct {
	entry reflect.Value
}

func (self *ReflectCallHandler) Call(ev *cellnet.Event) {

	// 这里的反射, 会影响非常少的效率, 但因为外部写法简单, 就算了
	self.entry.Call([]reflect.Value{reflect.ValueOf(ev.Msg)})

}

func NewReflectCallHandler(userCallback interface{}) cellnet.EventHandler {

	return &ReflectCallHandler{
		entry: reflect.ValueOf(userCallback),
	}

}
