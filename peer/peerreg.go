package peer

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"sort"
)

type PeerCreateFunc func() cellnet.Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

// 注册Peer创建器
func RegisterPeerCreator(f PeerCreateFunc) {

	// 临时实例化一个，获取类型
	dummyPeer := f()

	if _, ok := creatorByTypeName[dummyPeer.TypeName()]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[dummyPeer.TypeName()] = f
}

// Peer创建器列表
func PeerCreatorList() (ret []string) {

	for name := range creatorByTypeName {
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return
}

// 创建一个Peer
func NewPeer(peerType string) cellnet.Peer {
	peerCreator := creatorByTypeName[peerType]
	if peerCreator == nil {
		panic(fmt.Sprintf("Peer type not found, name: '%s'", peerType))
	}

	return peerCreator()
}

// 创建Peer后，设置基本属性
func NewGenericPeer(peerType, name, addr string, q cellnet.EventQueue) cellnet.GenericPeer {

	p := NewPeer(peerType)
	gp := p.(cellnet.GenericPeer)
	gp.SetName(name)
	gp.SetAddress(addr)
	gp.SetQueue(q)
	return gp
}
