package cellqueue

import (
	"fmt"
	"github.com/davyxu/x/container"
	"runtime/debug"
	"time"
)

type PanicNotifyFunc func(interface{}, *Queue)

type Queue struct {
	pipe *xcontainer.Pipe

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

// 开启事件循环
func (self *Queue) Run() *Queue {

	self.pipe.Run(func(raw interface{}) {

		switch value := raw.(type) {
		case func():
			if self.CapturePanic {
				defer func() {

					if err := recover(); err != nil {
						self.PanicNotify(err, self)
					}

				}()
			}

			value()
		case nil:
			break
		default:
			panic(fmt.Sprintf("unexpected type %T", value))
		}
	})

	return self
}

// 停止事件循环
func (self *Queue) Stop() *Queue {
	self.pipe.Stop()
	return self
}

// 创建默认长度的队列
func NewQueue() *Queue {

	return &Queue{
		pipe: xcontainer.NewPipe(),

		// 默认的崩溃捕获打印
		PanicNotify: func(raw interface{}, queue *Queue) {

			fmt.Printf("%s: %v \n%s\n", time.Now().Format("2006-01-02 15:04:05"), raw, string(debug.Stack()))
			debug.PrintStack()
		},
	}
}
