package test

import (
	"testing"
	"time"
)

type SignalTester struct {
	*testing.T

	signal chan int
}

func (self *SignalTester) WaitAndExpect(value int, msg string) bool {

	select {
	case v := <-self.signal:
		if v != value {
			self.Fail()
			self.Logf("%s\n", msg)
			return false
		}

	case <-time.After(2 * time.Second):
		self.Fail()
		return false
	}

	return true
}

func (self *SignalTester) Done(value int) {
	self.signal <- value
}

func NewSignalTester(t *testing.T) *SignalTester {

	return &SignalTester{
		T:      t,
		signal: make(chan int),
	}
}
