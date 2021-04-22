package tests

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

type SignalTester struct {
	*testing.T

	dataList []interface{}
	guard    sync.Mutex

	timeout time.Duration
}

func (self *SignalTester) SetTimeout(du time.Duration) {
	self.timeout = du
}

func (self *SignalTester) match(expect []interface{}) bool {
	self.guard.Lock()
	defer self.guard.Unlock()
	return len(self.dataList) == len(expect)
}

func (self *SignalTester) WaitAll(msg string, values ...interface{}) {

	timeoutTS := time.Now().Add(self.timeout)
	for {
		time.Sleep(time.Millisecond * 100)

		if self.timeout != 0 && time.Now().After(timeoutTS) {
			self.T.Errorf("%s timeout", msg)
			self.T.FailNow()
		}

		if !self.match(values) {
			continue
		}

		for index, v := range self.dataList {
			if !reflect.DeepEqual(v, values[index]) {
				self.T.Errorf("%s expect %v, got %v", msg, values[index], v)
				self.T.FailNow()
			}
		}

		break
	}
}

func (self *SignalTester) Done(value interface{}) {
	self.guard.Lock()
	self.dataList = append(self.dataList, value)
	self.guard.Unlock()
}

func NewSignalTester(t *testing.T) *SignalTester {

	return &SignalTester{
		T: t,
	}
}
