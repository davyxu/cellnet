package proc

import (
	"github.com/davyxu/cellnet"
)

type ProcessorBundleSetter interface {
	SetEventProcessor(v cellnet.MessageProcessor)
	SetEventHooker(v cellnet.EventHooker)
	SetEventHandler(v cellnet.EventHandler)
}

type ProcessorBinder func(initor ProcessorBundleSetter, userHandler cellnet.UserMessageHandler)

var (
	procByName = map[string]ProcessorBinder{}
)

func RegisterEventProcessor(procName string, f ProcessorBinder) {

	procByName[procName] = f
}

func BindProcessor(peer cellnet.Peer, procName string, userHandler cellnet.UserMessageHandler) {

	if proc, ok := procByName[procName]; ok {

		initor := peer.(ProcessorBundleSetter)

		proc(initor, userHandler)
	} else {
		panic("processor not found:" + procName)
	}
}
