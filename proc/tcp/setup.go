package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {

	transmitter := new(TCPMessageTransmitter)
	hooker := new(rpcEventHooker)

	proc.RegisterEventProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetEventTransmitter(transmitter)
		bundle.SetEventHooker(hooker)
		bundle.SetEventCallback(proc.NewQueuedEventCallback(userCallback))

	})
}
