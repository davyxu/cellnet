package cellnet

import (
	"time"
)

type Timer struct {
	exit chan bool
}

func (self *Timer) Stop() {
	self.exit <- true
}

func NewTimer(eq EventQueue, dur time.Duration, callback func(*Timer)) *Timer {

	self := &Timer{
        make(chan bool),
    }

	go func() {

		for {

			select {
			case <-time.After(dur):
				eq.Post(func() {

					callback(self)
				})
			case <-self.exit:
				goto exit
			}
		}
	exit:
	}()

	return self
}
