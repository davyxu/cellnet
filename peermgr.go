package cellnet

import (
	"fmt"
	"log"
)

var name2peerMap = make(map[string]Peer)

func NewPeer(evq *cellnet.EvQueue, name string, peerType string) Peer {

	if GetPeer(name) != nil {
		log.Println("duplicate peer name: ", name)
		return nil
	}

	p := createPeer(peerType, &peerProfile{name: name, queue: evq})

	if p == nil {
		log.Println("peer type not found: ", peerType)
		return nil
	}

	name2peerMap[name] = p

	return p
}

func GetPeer(name string) Peer {
	if v, ok := name2peerMap[name]; ok {
		return v
	}

	return nil
}

func StartAllPeer() {

}
