package cellnet

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// 事件队列
type EventQueue interface {
	// 事件队列开始工作
	StartLoop() EventQueue

	// 停止事件队列
	StopLoop() EventQueue

	// 等待退出
	Wait()

	// 投递事件, 通过队列到达消费者端
	Post(callback func())

	// 是否捕获异常
	EnableCapturePanic(v bool)
}

type CapturePanicNotifyFunc func(interface{}, EventQueue)

type eventQueue struct {
	*Pipe

	endSignal sync.WaitGroup

	capturePanic bool

	onPanic CapturePanicNotifyFunc
}

// 启动崩溃捕获
func (self *eventQueue) EnableCapturePanic(v bool) {
	self.capturePanic = v
}

// 设置捕获崩溃通知
func (self *eventQueue) SetCapturePanicNotify(callback CapturePanicNotifyFunc) {
	self.onPanic = callback
}

// 派发事件处理回调到队列中
func (self *eventQueue) Post(callback func()) {

	if callback == nil {
		return
	}

	self.Add(callback)
}

// 保护调用用户函数
func (self *eventQueue) protectedCall(callback func()) {

	if self.capturePanic {
		defer func() {

			if err := recover(); err != nil {
				self.onPanic(err, self)
			}

		}()
	}

	callback()
}

// 开启事件循环
func (self *eventQueue) StartLoop() EventQueue {

	self.endSignal.Add(1)

	go func() {

		var writeList []interface{}

		for {
			writeList = writeList[0:0]
			exit := self.Pick(&writeList)

			// 遍历要发送的数据
			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					self.protectedCall(t)
				case nil:
					break
				default:
					log.Printf("unexpected type %T", t)
				}
			}

			if exit {
				break
			}
		}

		self.endSignal.Done()
	}()

	return self
}

// 停止事件循环
func (self *eventQueue) StopLoop() EventQueue {
	self.Add(nil)
	return self
}

// 等待退出消息
func (self *eventQueue) Wait() {
	self.endSignal.Wait()
}

// 创建默认长度的队列
func NewEventQueue() EventQueue {

	return &eventQueue{
		Pipe: NewPipe(),

		// 默认的崩溃捕获打印
		onPanic: func(raw interface{}, queue EventQueue) {

			fmt.Printf("%s: %v \n%s\n", time.Now().Format("2006-01-02 15:04:05"), raw, string(debug.Stack()))
			debug.PrintStack()
		},
	}
}

// 在会话对应的Peer上的事件队列中执行callback，如果没有队列，则马上执行
func SessionQueuedCall(ses Session, callback func()) {
	if ses == nil {
		return
	}
	q := ses.Peer().(interface {
		Queue() EventQueue
	}).Queue()

	QueuedCall(q, callback)
}

// 有队列时队列调用，无队列时直接调用
func QueuedCall(queue EventQueue, callback func()) {
	if queue == nil {
		callback()
	} else {
		queue.Post(callback)
	}
}
