package timer

import (
	"github.com/davyxu/cellnet"
	"sync/atomic"
	"time"
)

type Loop struct {
	Context      interface{}
	Duration     time.Duration
	userCallback func(*Loop)

	running int64

	Queue cellnet.EventQueue
}

func (self *Loop) Running() bool {
	return atomic.LoadInt64(&self.running) != 0
}

func (self *Loop) setRunning(v bool) {

	if v {
		atomic.StoreInt64(&self.running, 1)
	} else {
		atomic.StoreInt64(&self.running, 0)
	}

}

func (self *Loop) Start() bool {

	if self.Running() {
		return false
	}

	atomic.StoreInt64(&self.running, 1)

	self.rawPost()

	return true
}

func (self *Loop) rawPost() {

	if self.Duration == 0 {
		panic("seconds can be zero in loop")
	}

	if self.Running() {
		After(self.Queue, self.Duration, func() {

			tick(self, false)
		}, nil)
	}
}

func (self *Loop) NextLoop() {

	self.Queue.Post(func() {
		tick(self, true)
	})
}

func (self *Loop) Stop() {

	self.setRunning(false)
}

func (self *Loop) Notify() *Loop {
	self.userCallback(self)
	return self
}

func tick(ctx interface{}, nextLoop bool) {

	loop := ctx.(*Loop)

	if !nextLoop && loop.Running() {

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
