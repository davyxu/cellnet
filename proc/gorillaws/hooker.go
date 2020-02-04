package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteRecvLogger("ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteSendLogger("ws", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}
