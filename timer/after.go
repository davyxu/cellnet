package celltimer

import (
	cellqueue "github.com/davyxu/cellnet/queue"
	"time"
)

type AfterStopper interface {
	Stop() bool
}

// 在给定的duration持续时间后, 执行callbackObj对象类型对应的函数回调
// q: 队列,在指定的队列goroutine执行, 空时,直接在当前goroutine
// context: 将context上下文传递到带有context指针的函数回调中
func After(q *cellqueue.Queue, duration time.Duration, callbackObj any, context any) AfterStopper {

	return time.AfterFunc(duration, func() {
		switch callback := callbackObj.(type) {
		case func():
			if callback != nil {
				cellqueue.QueuedCall(q, callback)
			}

		case func(any):
			if callback != nil {

				cellqueue.QueuedCall(q, func() {
					callback(context)
				})
			}
		default:
			panic("celltimer.After: require func() or func(any)")
		}
	})

}
