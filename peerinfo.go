package cellnet

import (
	"sync"
)

type PeerInfo interface {
	SetAddress(addr string)
	SetName(name string)
	SetQueue(q EventQueue)
}

type CorePeerInfo struct {
	name string

	address string

	tag      []tagData
	tagGuard sync.RWMutex

	queue EventQueue
}

func (self *CorePeerInfo) SetAddress(addr string) {
	self.address = addr
}

func (self *CorePeerInfo) Address() string {
	return self.address
}

// 获取通讯端的名称
func (self *CorePeerInfo) Name() string {
	return self.name
}

func (self *CorePeerInfo) SetName(name string) {
	self.name = name
}

// 获取队列
func (self *CorePeerInfo) EventQueue() EventQueue {
	return self.queue
}

func (self *CorePeerInfo) SetQueue(q EventQueue) {
	self.queue = q
}

func (self *CorePeerInfo) NameOrAddress() string {
	if self.name != "" {
		return self.name
	}

	return self.name
}
