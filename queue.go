package cellnet

import (
	"runtime/debug"
	"time"
)

type EventQueue interface {
	StartLoop()

	StopLoop(result int)

	// 等待退出
	Wait() int

	// 投递事件, 通过队列到达消费者端
	Post(callback func())

	// 延时投递
	DelayPost(dur time.Duration, callback func())
}

type queueData struct {
	data interface{}
}

type evQueue struct {
	queue chan func()

	exitSignal chan int

	capturePanic bool
}

// 派发到队列
func (self *evQueue) Post(callback func()) {

	self.queue <- callback
}

func (self *evQueue) DelayPost(dur time.Duration, callback func()) {
	go func() {

		time.AfterFunc(dur, func() {

			self.Post(callback)
		})

	}()
}

func (self *evQueue) protectedCall(callback func()) {

	if self.capturePanic {
		defer func() {

			if err := recover(); err != nil {
				//log.Fatalln(err)
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

func NewEventQueue() EventQueue {
	self := &evQueue{
		queue:      make(chan func(), 10),
		exitSignal: make(chan int),
	}

	return self

}
