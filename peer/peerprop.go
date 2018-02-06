package peer

import "github.com/davyxu/cellnet"

type CorePeerProperty struct {
	CorePropertySet
}

// 获取通讯端的名称
func (self *CorePeerProperty) Name() (ret string) {
	self.GetProperty("Name", &ret)
	return
}

// 获取队列
func (self *CorePeerProperty) Queue() (ret cellnet.EventQueue) {
	self.GetProperty("Queue", &ret)
	return
}
func (self *CorePeerProperty) Address() (ret string) {
	self.GetProperty("Address", &ret)
	return
}

func (self *CorePeerProperty) NameOrAddress() string {
	if name := self.Name(); name != "" {
		return name
	}

	return self.Address()
}
