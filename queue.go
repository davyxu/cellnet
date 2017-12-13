package cellnet

import (
	"runtime/debug"
	"sync"
)

type EventQueue interface {
	StartLoop()

	StopLoop(result int)

	// 等待退出
	Wait() int

	// 投递事件, 通过队列到达消费者端
	Post(callback func())

	// 是否捕获异常
	EnableCapturePanic(v bool)
}

type evQueue struct {
	queue chan func()

	endSignal sync.WaitGroup

	capturePanic bool

	result int
}

func (self *evQueue) EnableCapturePanic(v bool) {
	self.capturePanic = v
}

// 派发到队列
func (self *evQueue) Post(callback func()) {

	if callback == nil {
		return
	}

	self.queue <- callback
}

func (self *evQueue) protectedCall(callback func()) {

	if callback == nil {
		return
	}

	if self.capturePanic {
		defer func() {

			if err := recover(); err != nil {

				debug.PrintStack()
			}

		}()
	}

	callback()
}

func (self *evQueue) StartLoop() {

	self.endSignal.Add(1)

	go func() {
		for callback := range self.queue {

			if callback == nil {
				break
			}

			self.protectedCall(callback)
		}

		self.endSignal.Done()
	}()
}

func (self *evQueue) StopLoop(result int) {
	self.queue <- nil
	self.result = result
}

func (self *evQueue) Wait() int {
	self.endSignal.Wait()
	return self.result
}

const DefaultQueueSize = 100

func NewEventQueue() EventQueue {

	return NewEventQueueByLen(DefaultQueueSize)
}

func NewEventQueueByLen(l int) EventQueue {
	self := &evQueue{
		queue: make(chan func(), l),
	}

	return self
}
