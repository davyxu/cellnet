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

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", ":9001", nil).(cellnet.HTTPAcceptor)
	p.SetFileServe(".", ".")

	proc.BindProcessorHandler(p, "httpfile", nil)

	p.Start()
	queue.StartLoop()

	queue.Wait()
}
