package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {

	proc.RegisterEventProcessor("http", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetEventCallback(userCallback)
	})

}
