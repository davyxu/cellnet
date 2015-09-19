package cellnet

import (
	"github.com/davyxu/cellnet"
)

type Peer interface {
	Start(address string)
	Stop()
}

type PeerProfile struct {
	name  string
	queue *cellnet.EvQueue
}
