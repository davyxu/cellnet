package peer

import (
	"sync"
)

// 通信通讯端共享的数据
type CoreRunningTag struct {

	// 运行状态
	running      bool
	runningGuard sync.RWMutex

	// 停止过程同步
	stopping chan bool
}

func (self *CoreRunningTag) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *CoreRunningTag) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
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
	self.stopping = make(chan bool)
}

func (self *CoreRunningTag) EndStopping() {
	select {
	case self.stopping <- true:

	default:
		self.stopping = nil
	}
}
