package socket

import "github.com/davyxu/cellnet"

// 通讯端共享的数据
type socketPeer struct {
	cellnet.PeerConfig
	// 单独保存的保存cellnet.Peer接口
	peerInterface cellnet.Peer
}

// socket包内部派发事件
func (self *socketPeer) fireEvent(ev interface{}) interface{} {

	if self.Event == nil {
		return nil
	}

	return self.Event(ev)
}
