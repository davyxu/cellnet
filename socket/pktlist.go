package socket

import (
	"sync"

	"github.com/davyxu/cellnet"
)

type PacketList struct {
	list      []*cellnet.Packet
	listGuard sync.Mutex
	listCond  *sync.Cond
}

func (self *PacketList) Add(p *cellnet.Packet) {
	self.listGuard.Lock()
	self.list = append(self.list, p)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *PacketList) Reset() {
	self.list = self.list[0:0]
}

func (self *PacketList) BeginPick() []*cellnet.Packet {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	self.listGuard.Unlock()

	self.listGuard.Lock()

	return self.list
}

func (self *PacketList) EndPick() {

	self.Reset()
	self.listGuard.Unlock()
}

func NewPacketList() *PacketList {
	self := &PacketList{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}
