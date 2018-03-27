package peer

import (
	"github.com/davyxu/cellnet"
	"sort"
)

type PeerCreateFunc func() cellnet.Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

func RegisterPeerCreator(f PeerCreateFunc) {

	// 临时实例化一个，获取类型
	dummyPeer := f()

	if _, ok := creatorByTypeName[dummyPeer.TypeName()]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[dummyPeer.TypeName()] = f
}

func PeerCreatorList() (ret []string) {

	for name := range creatorByTypeName {
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return
}

func NewPeer(peerType string) cellnet.Peer {
	peerCreator := creatorByTypeName[peerType]
	if peerCreator == nil {
		panic("Peer name not found: " + peerType)
	}

	return peerCreator()
}

func NewGenericPeer(peerType, name, addr string, q cellnet.EventQueue) cellnet.GenericPeer {

	p := NewPeer(peerType)
	gp := p.(cellnet.GenericPeer)
	gp.SetName(name)
	gp.SetAddress(addr)
	gp.SetQueue(q)
	return gp
}
