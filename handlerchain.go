package cellnet

import "bytes"

type HandlerChain struct {
	list []EventHandler
}

func (self *HandlerChain) Add(h EventHandler) {
	self.list = append(self.list, h)
}

func (self *HandlerChain) String() string {

	var buff bytes.Buffer

	buff.WriteString("	")

	for index, h := range self.list {

		if index > 0 {
			buff.WriteString(" -> ")
		}

		buff.WriteString(HandlerString(h))
	}

	return buff.String()
}

func (self *HandlerChain) Call(ev *Event) {

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

func NewHandlerChain(h ...EventHandler) *HandlerChain {
	return &HandlerChain{
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
