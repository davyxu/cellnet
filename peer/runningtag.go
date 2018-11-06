package peer

import (
	"sync"
	"sync/atomic"
)

// 通信通讯端共享的数据
type CoreRunningTag struct {

	// 运行状态
	running int64

	stoppingWaitor sync.WaitGroup
	stopping       int64
}

func (self *CoreRunningTag) IsRunning() bool {

	return atomic.LoadInt64(&self.running) != 0
}

func (self *CoreRunningTag) SetRunning(v bool) {

	if v {
		atomic.StoreInt64(&self.running, 1)
	} else {
		atomic.StoreInt64(&self.running, 0)
	}

}

func (self *CoreRunningTag) WaitStopFinished() {
	// 如果正在停止时, 等待停止完成
	self.stoppingWaitor.Wait()
}

func (self *CoreRunningTag) IsStopping() bool {
	return atomic.LoadInt64(&self.stopping) != 0
}

func (self *CoreRunningTag) StartStopping() {
	self.stoppingWaitor.Add(1)
	atomic.StoreInt64(&self.stopping, 1)
}

func (self *CoreRunningTag) EndStopping() {

	if self.IsStopping() {
		self.stoppingWaitor.Done()
		atomic.StoreInt64(&self.stopping, 0)
	}

}
