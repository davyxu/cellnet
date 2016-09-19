package benchmark

import (
	"sync"
	"time"

	"github.com/davyxu/cellnet"
)

type QPSMeter struct {
	qpsGuard sync.Mutex
	qps      int
	total    int

	count int
}

func (self *QPSMeter) Acc() int {
	self.qpsGuard.Lock()

	defer self.qpsGuard.Unlock()

	self.qps++

	return self.count
}

// 一轮计算
func (self *QPSMeter) Turn() (ret int) {
	self.qpsGuard.Lock()

	if self.qps > 0 {
		ret = self.qps
	}

	self.total += self.qps

	self.qps = 0
	self.count++

	self.qpsGuard.Unlock()

	return
}

// 均值
func (self *QPSMeter) Average() int {

	self.qpsGuard.Lock()

	defer self.qpsGuard.Unlock()

	if self.count == 0 {
		return 0
	}

	return self.total / self.count
}

func NewQPSMeter(pipe cellnet.EventPipe, callback func(int)) *QPSMeter {

	self := &QPSMeter{}

	timeEvq := pipe.AddQueue()

	cellnet.NewTimer(timeEvq, time.Second, func(t *cellnet.Timer) {

		qps := self.Turn()

		callback(qps)

	})

	return self
}
