package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/binary/coredef"
)

type RPCMatchHandler struct {
	*cellnet.DispatcherHandler
}

func (self *RPCMatchHandler) Call(ev *cellnet.SessionEvent) {

	msg := ev.Msg.(*coredef.RemoteCallACK)

	// rpc需要匹配
	h := self.DispatcherHandler.GetHandlerByIndex(int(msg.MsgID), int(msg.CallID))
	if h != nil {
		cellnet.HandlerChainCall(h, ev)
	}
}

func NewRPCMatchHandler() cellnet.EventHandler {

	return &RPCMatchHandler{
		DispatcherHandler: cellnet.NewDispatcherHandler(),
	}

}
