package socket

import (
	"sync"

	"github.com/davyxu/cellnet"
)

type eventList struct {
	list      []*cellnet.Event
	listGuard sync.Mutex
	listCond  *sync.Cond
}

func (self *eventList) Add(ev *cellnet.Event) {
	self.listGuard.Lock()
	self.list = append(self.list, ev)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *eventList) Reset() {
	self.list = self.list[0:0]
}

func (self *eventList) BeginPick() []*cellnet.Event {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	self.listGuard.Unlock()

	self.listGuard.Lock()

	return self.list
}

func (self *eventList) EndPick() {

	self.Reset()
	self.listGuard.Unlock()
}

func NewPacketList() *eventList {
	self := &eventList{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}
