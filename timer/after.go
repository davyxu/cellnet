package timer

import (
	"github.com/davyxu/cellnet"
	"time"
)

type AfterStopper interface {
	Stop() bool
}

func After(q cellnet.EventQueue, duration time.Duration, callback func()) AfterStopper {

	if q == nil {
		return nil
	}

	afterTimer := time.NewTimer(duration)

	go func() {

		<-afterTimer.C

		q.Post(callback)

	}()

	return afterTimer

}
