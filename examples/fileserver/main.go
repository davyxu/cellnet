package main

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/httpfilepeer"
)

func main() {
	queue := cellnet.NewEventQueue()
	peer := cellnet.NewPeer("http.file.Acceptor")
	peerInfo := peer.(cellnet.PeerInfo)
	peerInfo.SetName("httpfile")
	peerInfo.SetAddress(":9001")
	peer.Start()
	queue.StartLoop()
	queue.Wait()
}
