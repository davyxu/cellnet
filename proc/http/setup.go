package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

func init() {

	proc.RegisterEventProcessor("http", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {
		initor.SetEventHandler(userHandler)
	})

}
