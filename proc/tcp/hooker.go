package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc/rpc"
)

// 自动区分rpc日志和tcp日志
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	if !rpc.ProcInboundEvent(inputEvent) {
		msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	}

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	if !rpc.ProcOutboundEvent(inputEvent) {
		msglog.WriteSendLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	}

	return inputEvent
}
