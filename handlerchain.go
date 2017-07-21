package cellnet

import (
	"bytes"
	"fmt"
	"reflect"
)

type HandlerChain struct {
	id   int64
	list []EventHandler
}

// 添加1个
func (self *HandlerChain) Add(h EventHandler) {
	self.list = append(self.list, h)
}

// 添加多个
func (self *HandlerChain) AddBatch(h ...EventHandler) {
	self.list = append(self.list, h...)
}

// 启动匹配类型
func (self *HandlerChain) AddAny(objlist ...interface{}) {

	for _, obj := range objlist {

		switch v := obj.(type) {
		case EventHandler:
			self.Add(v)
		case []EventHandler:
			self.list = append(self.list, v...)
		default:
			panic("unknown hander chain input type: " + reflect.ValueOf(v).String())
		}

	}

}

func (self *HandlerChain) String() string {

	var buff bytes.Buffer

	buff.WriteString(fmt.Sprintf("	 chain: %d ", self.id))

	for index, h := range self.list {

		if index > 0 {
			buff.WriteString(" -> ")
		}

		buff.WriteString(HandlerString(h))
	}

	return buff.String()
}

func (self *HandlerChain) Call(ev *Event) {

	ev.chainid = self.id

	for _, h := range self.list {

		HandlerLog(h, ev)

		h.Call(ev)

		if ev.Result() == Result_NextChain {
			ev.SetResult(Result_OK)
			break
		}

		if ev.Result() != Result_OK {
			break
		}
	}

}

var chainidgen int64 = 500

func genChainID() int64 {
	chainidgen++
	return chainidgen
}

func NewHandlerChain(objlist ...interface{}) *HandlerChain {
	self := &HandlerChain{
		id: genChainID(),
	}

	self.AddAny(objlist...)

	return self
}

type HandlerChainList []*HandlerChain

func (self HandlerChainList) Call(ev *Event) {

	for _, chain := range self {

		cloned := ev.Clone()

		chain.Call(cloned)
	}

}

func (self HandlerChainList) String() string {

	var buff bytes.Buffer

	for _, chain := range self {

		buff.WriteString(chain.String())

		buff.WriteString("\n")
	}

	return buff.String()
}
