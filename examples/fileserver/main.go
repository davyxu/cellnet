package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
)

func main() {
	queue := cellnet.NewEventQueue()
	peerIns := peer.NewPeer("http.Acceptor")
	pset := peerIns.(cellnet.PropertySet)
	pset.SetProperty("Address", ":9001")
	pset.SetProperty("Name", "httpfile")
	pset.SetProperty("HttpDir", ".")
	pset.SetProperty("HttpRoot", ".")

	proc.BindProcessor(peerIns, "httpfile", nil)

	peerIns.Start()
	queue.StartLoop()

	queue.Wait()
}
