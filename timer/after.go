package timer

import (
	"github.com/davyxu/cellnet"
	"time"
)

type AfterStopper interface {
	Stop() bool
}

// 在给定的duration持续时间后, 执行callbackObj对象类型对应的函数回调
// q: 队列,在指定的队列goroutine执行, 空时,直接在当前goroutine
// context: 将context上下文传递到带有context指针的函数回调中
func After(q cellnet.EventQueue, duration time.Duration, callbackObj interface{}, context interface{}) AfterStopper {

	return time.AfterFunc(duration, func() {
		switch callback := callbackObj.(type) {
		case func():
			if callback != nil {
				cellnet.QueuedCall(q, callback)
			}

		case func(interface{}):
			if callback != nil {

				cellnet.QueuedCall(q, func() {
					callback(context)
				})
			}
		default:
			panic("timer.After: require func() or func(interface{})")
		}
	})

}
