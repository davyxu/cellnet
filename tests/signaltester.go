package tests

import (
	"testing"
	"time"
)

type SignalTester struct {
	*testing.T

	signal chan interface{}

	timeout time.Duration
}

func (self *SignalTester) SetTimeout(du time.Duration) {
	self.timeout = du
}

func (self *SignalTester) WaitAndExpect(msg string, values ...interface{}) bool {

	var recvValues = map[interface{}]bool{}

	for _, value := range values {
		select {
		case v := <-self.signal:
			recvValues[v] = true

		case <-time.After(self.timeout):
			self.Errorf("signal timeout: %d %s", value, msg)
			self.Fail()
			return false
		}
	}

	for _, value := range values {

		if _, ok := recvValues[value]; !ok {
			self.FailNow()
		}

	}

	return true
}

func (self *SignalTester) Done(value interface{}) {
	self.signal <- value
}

func NewSignalTester(t *testing.T) *SignalTester {

	return &SignalTester{
		T:       t,
		timeout: 2 * time.Second,
		signal:  make(chan interface{}),
	}
}
