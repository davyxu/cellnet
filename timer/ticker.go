package timer

import (
	"github.com/davyxu/cellnet"
	"time"
)

type TickStopper interface {
	Stop()
}

func Tick(q cellnet.EventQueue, duration time.Duration, callbackObj interface{}) TickStopper {

	if q == nil {
		return nil
	}

	ticker := time.NewTicker(duration)

	switch callback := callbackObj.(type) {
	case func():
		go func() {

			for {
				select {
				case <-ticker.C:
					q.Post(callback)
				}
			}

		}()

	case func(TickStopper):
		go func() {

			for {
				select {
				case <-ticker.C:
					q.Post(func() {
						callback(ticker)
					})
				}
			}

		}()
	default:
		panic("Require func() or func(TickStopper) for callback")
	}

	return ticker
}
