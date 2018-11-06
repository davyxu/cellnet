package peer

import "github.com/davyxu/cellnet"

type CorePeerProperty struct {
	name  string
	queue cellnet.EventQueue
	addr  string
}

// 获取通讯端的名称
func (self *CorePeerProperty) Name() string {
	return self.name
}

// 获取队列
func (self *CorePeerProperty) Queue() cellnet.EventQueue {
	return self.queue
}

// 获取SetAddress中的侦听或者连接地址
func (self *CorePeerProperty) Address() string {

	return self.addr
}

func (self *CorePeerProperty) SetName(v string) {
	self.name = v
}

func (self *CorePeerProperty) SetQueue(v cellnet.EventQueue) {
	self.queue = v
}

func (self *CorePeerProperty) SetAddress(v string) {
	self.addr = v
}
