package rpc

import "github.com/davyxu/cellnet"

type RetChanHandler struct {
	cellnet.BaseEventHandler
	ret chan interface{}
}

func (self *RetChanHandler) Call(ev *cellnet.SessionEvent) error {

	self.ret <- ev.Msg

	return self.CallNext(ev)
}

func NewRetChanHandler(ret chan interface{}) cellnet.EventHandler {
	return &RetChanHandler{
		ret: ret,
	}

}
