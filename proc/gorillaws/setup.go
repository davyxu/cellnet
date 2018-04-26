package gorillaws

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {

	transmitter := new(WSMessageTransmitter)
	hooker := new(MsgHooker)

	proc.RegisterProcessor("gorillaws.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(hooker)
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
