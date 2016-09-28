package socket

import (
	"github.com/davyxu/cellnet"
)

// Peer间的共享数据
type peerProfile struct {
	cellnet.EventQueue // 实现事件注册和注入
	name               string
	maxPacketSize      int
}

func (self *peerProfile) SetName(name string) {
	self.name = name
}

func (self *peerProfile) Name() string {
	return self.name
}

func (self *peerProfile) SetMaxPacketSize(size int) {
	self.maxPacketSize = size
}

func (self *peerProfile) MaxPacketSize() int {
	return self.maxPacketSize
}

func newPeerProfile(queue cellnet.EventQueue) *peerProfile {

	self := &peerProfile{EventQueue: queue}

	return self
}
