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

type lineralTask struct {
	q *evQueue
	e interface{}
}

type lineraPipe struct {
	exitSignal chan int

	dataChan chan *lineralTask

	capturePanic bool
}

func (self *lineraPipe) AddQueue() EventQueue {

	q := newEventQueue()

	go func(q *evQueue) {
		for v := range q.queue {
			self.dataChan <- &lineralTask{q: q, e: v}
		}
	}(q)

	return q
}

func (self *lineraPipe) EnableCaputrePanic(enable bool) {
	self.capturePanic = enable
}

func (self *lineraPipe) Start() {

	go func() {
		for v := range self.dataChan {
			self.protectedCall(v.q, v.e)
		}
	}()
}
func (self *lineraPipe) protectedCall(q *evQueue, data interface{}) {
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

func (self *lineraPipe) Stop(result int) {
	self.exitSignal <- result
}

func (self *lineraPipe) Wait() int {
	return <-self.exitSignal
}

func NewEventPipe() EventPipe {
	return &lineraPipe{
		exitSignal: make(chan int),
		dataChan:   make(chan *lineralTask),
	}
}
