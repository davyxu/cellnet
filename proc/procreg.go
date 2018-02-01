package proc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

type DuplexEventProcessor func(cellnet.EventProc, cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc)

var evtprocByName = map[string]DuplexEventProcessor{}

func RegisterEventProcessor(name string, f DuplexEventProcessor) {

	if peer.PeerCreatorExists(name) {
		panic("Duplicate peer type")
	}

	evtprocByName[name] = f
}

func MakeEventProcessor(name string, inbound, outbound cellnet.EventProc) (cellnet.EventProc, cellnet.EventProc) {

	f := evtprocByName[name]
	if f == nil {
		panic("Event processor not found: " + name)
	}

	return f(inbound, outbound)
}
