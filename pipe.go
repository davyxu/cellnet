package cellnet

import (
	"runtime/debug"
)

type EventPipe interface {
	AddQueue() EventQueue

	Start()

	Stop(int)

	Wait() int

	// 开启捕获错误, 错误不会崩溃
	EnableCaputrePanic(enable bool)
}

type linearTask struct {
	q *evQueue
	e interface{}
}

type linearPipe struct {
	exitSignal chan int

	dataChan chan *linearTask

	capturePanic bool
}

func (self *linearPipe) AddQueue() EventQueue {

	q := newEventQueue()

	go func(q *evQueue) {
		for v := range q.queue {
			self.dataChan <- &linearTask{q: q, e: v}
		}
	}(q)

	return q
}

func (self *linearPipe) EnableCaputrePanic(enable bool) {
	self.capturePanic = enable
}

func (self *linearPipe) Start() {

	go func() {
		for v := range self.dataChan {
			self.protectedCall(v.q, v.e)
		}
	}()
}
func (self *linearPipe) protectedCall(q *evQueue, data interface{}) {
	if self.capturePanic {
		defer func() {

			if err := recover(); err != nil {
				log.Fatalln(err)
				debug.PrintStack()
			}

		}()
	}

	q.CallData(data)
}

func (self *linearPipe) Stop(result int) {
	self.exitSignal <- result
}

func (self *linearPipe) Wait() int {
	return <-self.exitSignal
}

func NewEventPipe() EventPipe {
	return &linearPipe{
		exitSignal: make(chan int),
		dataChan:   make(chan *linearTask),
	}
}
