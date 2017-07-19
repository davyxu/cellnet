package cellnet

import (
	"bytes"
	"fmt"
)

type HandlerChain struct {
	id   int64
	list []EventHandler
}

func (self *HandlerChain) Add(h EventHandler) {
	self.list = append(self.list, h)
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

func NewHandlerChain(h ...EventHandler) *HandlerChain {
	return &HandlerChain{
		id:   genChainID(),
		list: h,
	}
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
