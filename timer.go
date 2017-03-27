package cellnet

import (
	"time"
)

type Timer struct {
	tick *time.Ticker
	done chan struct{}
}

func (self *Timer) Stop() {
	self.done <- struct{}{}
}

func NewTimer(eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {

	self := &Timer{
		tick: time.NewTicker(dur),
		done: make(chan struct{}),
	}

	go func() {
		defer self.tick.Stop()
		for {
			select {
			case <-self.tick.C:
				eq.Post(nil, func() {
					callback(self)
				})
			case <-self.done:
				return
			}
		}
	}()

	return self
}
