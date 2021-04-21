package cellqueue

import (
	"fmt"
	"github.com/davyxu/x/frame"
	"runtime/debug"
	"sync"
	"time"
)

type PanicNotifyFunc func(interface{}, *Queue)

type Queue struct {
	pipe      *xframe.Pipe
	endSignal sync.WaitGroup

	// 启动崩溃捕获
	CapturePanic bool

	// 设置捕获崩溃通知
	PanicNotify PanicNotifyFunc
}

// 派发事件处理回调到队列中
func (self *Queue) Post(callback func()) {

	if callback == nil {
		return
	}

	self.pipe.Add(callback)
}

// 保护调用用户函数
func (self *Queue) protectedCall(callback func()) {

	if self.CapturePanic {
		defer func() {

			if err := recover(); err != nil {
				self.PanicNotify(err, self)
			}

		}()
	}

	callback()
}

// 开启事件循环
func (self *Queue) Run() *Queue {

	self.endSignal.Add(1)

	go func() {

		var writeList []interface{}

		for {
			writeList = writeList[0:0]
			exit := self.pipe.Pick(&writeList)

			// 遍历要发送的数据
			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					self.protectedCall(t)
				case nil:
					break
				default:
					panic(fmt.Sprintf("unexpected type %T", t))
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
func (self *Queue) Stop() *Queue {
	self.pipe.Add(nil)
	return self
}

// 等待退出消息
func (self *Queue) Wait() {
	self.endSignal.Wait()
}

// 创建默认长度的队列
func NewQueue() *Queue {

	return &Queue{
		pipe: xframe.NewPipe(),

		// 默认的崩溃捕获打印
		PanicNotify: func(raw interface{}, queue *Queue) {

			fmt.Printf("%s: %v \n%s\n", time.Now().Format("2006-01-02 15:04:05"), raw, string(debug.Stack()))
			debug.PrintStack()
		},
	}
}
