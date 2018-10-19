package cellnet

import (
	"math"
	"runtime"
	"sync/atomic"
)

type EventQueueGroup interface {
	// 事件队列开始工作
	StartLoop() EventQueueGroup

	// 停止事件队列
	StopLoop() EventQueueGroup

	// 等待退出
	Wait()

	GetQueue() EventQueue
}

type eventQueueGroup struct {
	list    []EventQueue
	size    int32
	counter int32
}

func NewEventGroup(size int) EventQueueGroup {
	if size < 1 {
		size = runtime.NumCPU()
	}
	group := &eventQueueGroup{}
	group.size = int32(size)
	group.list = make([]EventQueue, size)
	for i := int32(0); i < group.size; i++ {
		group.list[i] = NewEventQueue()
	}
	return group
}

func (self *eventQueueGroup) StartLoop() EventQueueGroup {
	if self.size <= 0 {
		return self
	}
	for i := int32(0); i < self.size; i++ {
		self.list[i].StartLoop()
	}
	return self
}

func (self *eventQueueGroup) StopLoop() EventQueueGroup {
	if self.size <= 0 {
		return self
	}
	for i := int32(0); i < self.size; i++ {
		self.list[i].StopLoop()
	}
	return self
}

//如果没有成员 无法阻塞
func (self *eventQueueGroup) Wait() {
	if self.size <= 0 {
		return
	}
	self.list[0].Wait()
}

func (self *eventQueueGroup) GetQueue() EventQueue {
	oldCounter := self.counter
	newCounter := int32(math.Min(float64(self.size), float64(oldCounter+1)))
	for {
		if atomic.CompareAndSwapInt32(&self.counter, oldCounter, newCounter) {
			break
		}
		oldCounter = self.counter
		newCounter = int32(math.Min(float64(self.size), float64(oldCounter+1)))
	}

	return self.list[int(newCounter)]
}
