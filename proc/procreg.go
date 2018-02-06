package proc

import (
	"github.com/davyxu/cellnet"
)

type ProcessorBundleInitor interface {
	SetEventProcessor(v cellnet.MessageProcessor)
	SetEventHooker(v cellnet.EventHooker)
	SetEventHandler(v cellnet.EventHandler)
}

type ProcessorBinder func(initor ProcessorBundleInitor, userHandler cellnet.UserMessageHandler)

var (
	procByName = map[string]ProcessorBinder{}
)

func RegisterEventProcessor(procName string, f ProcessorBinder) {

	procByName[procName] = f
}

func BindProcessor(peer cellnet.Peer, procName string, userHandler cellnet.UserMessageHandler) {

	if proc, ok := procByName[procName]; ok {

		initor := peer.(ProcessorBundleInitor)

		proc(initor, userHandler)
	} else {
		panic("processor not found:" + procName)
	}
}
