package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/relay"
	"github.com/davyxu/cellnet/rpc"
)

// 带有RPC和relay功能
type MsgHooker struct {
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	var handled bool

	if inputEvent, handled = rpc.ResolveInboundEvent(inputEvent); !handled {

		if inputEvent, handled = relay.ResoleveInboundEvent(inputEvent); !handled {
			msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
		}
	}

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	if !rpc.ResolveOutboundEvent(inputEvent) {

		if !relay.ResolveOutboundEvent(inputEvent) {
			msglog.WriteSendLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
		}
	}

	return inputEvent
}
