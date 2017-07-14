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

func (self *eventList) Pick() (ret []*cellnet.Event, exit bool) {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	self.listGuard.Unlock()

	self.listGuard.Lock()

	// 复制出队列

	for _, ev := range self.list {

		if ev == nil {
			exit = true
			break
		} else {
			ret = append(ret, ev)
		}
	}

	self.Reset()
	self.listGuard.Unlock()

	return
}

func NewPacketList() *eventList {
	self := &eventList{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}
