package timer

import (
	"github.com/davyxu/cellnet"
	"time"
)

type Loop struct {
	Context      interface{}
	Duration     time.Duration
	userCallback func(*Loop)

	running bool

	Queue cellnet.EventQueue
}

func (self *Loop) Running() bool {
	return self.running
}

func (self *Loop) Start() bool {

	if self.running {
		return false
	}

	self.running = true

	self.rawPost()

	return true
}

func (self *Loop) rawPost() {

	if self.Duration == 0 {
		panic("seconds can be zero in loop")
	}

	After(self.Queue, self.Duration, func() {

		tick(self, false)
	}, nil)
}

func (self *Loop) NextLoop() {

	self.Queue.Post(func() {
		tick(self, true)
	})
}

func (self *Loop) Stop() {

	self.running = false
}

func (self *Loop) Notify() *Loop {
	self.userCallback(self)
	return self
}

func tick(ctx interface{}, nextLoop bool) {

	loop := ctx.(*Loop)

	if !nextLoop && loop.running {

		// 即便在Notify中发生了崩溃，也会使用defer再次继续循环
		defer loop.rawPost()
	}

	loop.Notify()
}

func NewLoop(q cellnet.EventQueue, duration time.Duration, callback func(*Loop), context interface{}) *Loop {

	self := &Loop{
		Context:      context,
		Duration:     duration,
		userCallback: callback,
		Queue:        q,
	}

	return self
}
