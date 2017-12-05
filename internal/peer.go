package internal

import "github.com/davyxu/cellnet"

// 通讯端共享的数据
type PeerShare struct {
	cellnet.PeerConfig
	SessionManager

	// 单独保存的保存cellnet.Peer接口
	peerInterface cellnet.Peer
}

// socket包内部派发事件
func (self *PeerShare) FireEvent(ev interface{}) interface{} {

	if self.Event == nil {
		return nil
	}

	return self.Event(ev)
}

func (self *PeerShare) Init(p cellnet.Peer, config cellnet.PeerConfig) {
	self.SessionManager = NewSessionManager()
	self.peerInterface = p
	self.PeerConfig = config
}
