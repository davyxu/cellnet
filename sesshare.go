package cellnet

type CoreSessionShare struct {
	CoreTagger
	id int64
	// 归属的通讯端
	PeerShare *CoreCommunicatePeer
}

// 取会话归属的通讯端
func (self *CoreSessionShare) Peer() Peer {
	return self.PeerShare.Peer()
}

func (self *CoreSessionShare) ID() int64 {
	return self.id
}

func (self *CoreSessionShare) SetID(id int64) {
	self.id = id
}
