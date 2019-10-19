package timer

import (
	"github.com/davyxu/cellnet"
	"sync/atomic"
	"time"
)

// 轻量级的持续Tick循环
type Loop struct {
	Context        interface{}
	Duration       time.Duration
	notifyCallback func(*Loop)

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

// 开始Tick
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

func (self *Loop) Resume() {

	self.setRunning(true)
}

// 马上调用一次用户回调
func (self *Loop) Notify() *Loop {
	self.notifyCallback(self)
	return self
}

func (self *Loop) SetNotifyFunc(notifyCallback func(*Loop)) *Loop {
	self.notifyCallback = notifyCallback
	return self
}

func (self *Loop) NotifyFunc() func(*Loop) {
	return self.notifyCallback
}

func tick(ctx interface{}, nextLoop bool) {

	loop := ctx.(*Loop)

	if !nextLoop && loop.Running() {

		// 即便在Notify中发生了崩溃，也会使用defer再次继续循环
		defer loop.rawPost()
	}

	loop.Notify()
}

// 执行一个循环, 持续调用callback, 周期是duration
// context: 将context上下文传递到带有context指针的函数回调中
func NewLoop(q cellnet.EventQueue, duration time.Duration, notifyCallback func(*Loop), context interface{}) *Loop {

	self := &Loop{
		Context:        context,
		Duration:       duration,
		notifyCallback: notifyCallback,
		Queue:          q,
	}

	return self
}
