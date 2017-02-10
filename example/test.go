package test

import (
	"testing"
	"time"
)

type SignalTester struct {
	*testing.T

	signal chan int

	timeout time.Duration
}

func (self *SignalTester) SetTimeout(du time.Duration) {
	self.timeout = du
}

func (self *SignalTester) WaitAndExpect(value int, msg string) bool {

	select {
	case v := <-self.signal:
		if v != value {
			self.Fail()
			self.Logf("%s\n", msg)
			return false
		}

	case <-time.After(self.timeout):
		self.Logf("signal timeout: %d %s", value, msg)
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
		T:       t,
		timeout: 2 * time.Second,
		signal:  make(chan int),
	}
}
