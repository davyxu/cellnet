package internal

import "github.com/davyxu/cellnet"

type PeerInfo struct {
	name string

	address string

	tag interface{}

	queue cellnet.EventQueue
}

func (self *PeerInfo) SetAddress(addr string) {
	self.address = addr
}

func (self *PeerInfo) Address() string {
	return self.address
}

// 获取通讯端的名称
func (self *PeerInfo) Name() string {
	return self.name
}

func (self *PeerInfo) SetName(name string) {
	self.name = name
}

// 获取队列
func (self *PeerInfo) EventQueue() cellnet.EventQueue {
	return self.queue
}

func (self *PeerInfo) SetQueue(q cellnet.EventQueue) {
	self.queue = q
}

func (self *PeerInfo) SetTag(tag interface{}) {
	self.tag = tag
}

func (self *PeerInfo) Tag() interface{} {
	return self.tag
}

func (self *PeerInfo) NameOrAddress() string {
	if self.name != "" {
		return self.name
	}

	return self.name
}
