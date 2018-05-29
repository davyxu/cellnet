package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
)

var shareDir = flag.String("share", ".", "folder to share")
var port = flag.Int("port", 9091, "listen port")

func main() {

	flag.Parse()

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", fmt.Sprintf(":%d", *port), nil).(cellnet.HTTPAcceptor)
	p.SetFileServe(".", *shareDir)

	proc.BindProcessorHandler(p, "http", nil)

	p.Start()
	queue.StartLoop()

	queue.Wait()
}
