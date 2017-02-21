package rpc

import "github.com/davyxu/cellnet"

type RetChanHandler struct {
	cellnet.BaseEventHandler
	ret chan interface{}
}

func (self *RetChanHandler) Call(ev *cellnet.SessionEvent) {

	self.ret <- ev.Msg
}

func NewRetChanHandler(ret chan interface{}) cellnet.EventHandler {
	return &RetChanHandler{
		ret: ret,
	}

}
