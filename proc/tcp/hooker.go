package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc/rpc"
)

type rpcEventHooker struct {
	rpc.RPCHooker
}

func (self rpcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())

	self.RPCHooker.OnInboundEvent(inputEvent)

	return inputEvent
}

func (self rpcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteSendLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())

	self.RPCHooker.OnOutboundEvent(inputEvent)

	return inputEvent
}
