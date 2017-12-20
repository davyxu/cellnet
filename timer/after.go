package timer

import (
	"github.com/davyxu/cellnet"
	"time"
)

type AfterStopper interface {
	Stop() bool
}

func After(q cellnet.EventQueue, duration time.Duration, callbackObj interface{}, context interface{}) AfterStopper {

	afterTimer := time.NewTimer(duration)

	go func() {

		<-afterTimer.C

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

	}()

	return afterTimer

}
