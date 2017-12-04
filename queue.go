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

	// 是否捕获异常
	EnableCapturePanic(v bool)
}

type eventQueue struct {
	queue chan func()

	exitSignal chan int

	capturePanic bool
}

// 启动崩溃捕获
func (q *eventQueue) EnableCapturePanic(v bool) {
	q.capturePanic = v
}

// 派发事件处理回调到队列中
func (q *eventQueue) Post(callback func()) {

	if callback == nil {
		return
	}

	q.queue <- callback
}

// 保护调用用户函数
func (q *eventQueue) protectedCall(callback func()) {

	if callback == nil {
		return
	}

	if q.capturePanic {
		defer func() {

			if err := recover(); err != nil {

				debug.PrintStack()
			}

		}()
	}

	callback()
}

// 开启事件循环
func (q *eventQueue) StartLoop() {

	go func() {
		for callback := range q.queue {
			q.protectedCall(callback)
		}
	}()
}

// 停止事件循环
func (q *eventQueue) StopLoop(result int) {
	q.exitSignal <- result
}

// 等待退出消息
func (q *eventQueue) Wait() int {
	return <-q.exitSignal
}

const DefaultQueueSize = 100

// 创建默认长度的队列
func NewEventQueue() EventQueue {

	return NewEventQueueByLen(DefaultQueueSize)
}

// 创建指定长度的队列
func NewEventQueueByLen(l int) EventQueue {
	self := &eventQueue{
		queue:      make(chan func(), l),
		exitSignal: make(chan int),
	}

	return self
}
