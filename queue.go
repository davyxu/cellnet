package cellnet

import (
	"runtime/debug"
)

type EventQueue interface {
	StartLoop()

	StopLoop(result int)

	// 等待退出
	Wait() int

	// 投递事件, 通过队列到达消费者端
	Post(callback func())
}

type evQueue struct {
	queue chan func()

	exitSignal chan int

	capturePanic bool
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

	go func() {
		for callback := range self.queue {
			self.protectedCall(callback)
		}
	}()
}

func (self *evQueue) StopLoop(result int) {
	self.exitSignal <- result
}

func (self *evQueue) Wait() int {
	return <-self.exitSignal
}

const DefaultQueueSize = 100

func NewEventQueue() EventQueue {

	return NewEventQueueByLen(DefaultQueueSize)
}

func NewEventQueueByLen(l int) EventQueue {
	self := &evQueue{
		queue:      make(chan func(), l),
		exitSignal: make(chan int),
	}

	return self
}
