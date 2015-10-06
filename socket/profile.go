package socket

import (
	"github.com/davyxu/cellnet"
)

// Peer间的共享数据
type peerProfile struct {
	cellnet.EventQueue // 实现事件注册和注入
	name               string
}

func (self *peerProfile) SetName(name string) {
	self.name = name
}

func (self *peerProfile) Name() string {
	return self.name
}

func newPeerProfile(queue cellnet.EventQueue) *peerProfile {

	return &peerProfile{EventQueue: queue}
}
