package cellnet

var type2peercreator = make(map[string]func(*PeerProfile) Peer)

// 为新的peertype注册
func RegisterPeerType(typeName string, createFunc func(*PeerProfile) Peer) {

	type2peercreator[typeName] = createFunc
}

func createPeer(typeName string, pf *PeerProfile) Peer {

	if creator, ok := type2peercreator[typeName]; ok {
		return creator(pf)
	}

	return nil
}
