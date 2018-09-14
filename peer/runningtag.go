package peer

import (
	"sync/atomic"
)

// 通信通讯端共享的数据
type CoreRunningTag struct {

	// 运行状态
	running int64

	// 停止过程同步
	stopping chan struct{}
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
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}

func (self *CoreRunningTag) IsStopping() bool {
	return self.stopping != nil
}

func (self *CoreRunningTag) StartStopping() {
	self.stopping = make(chan struct{})
}

func (self *CoreRunningTag) EndStopping() {
	select {
	case self.stopping <- struct{}{}:

	default:
		self.stopping = nil
	}
}
