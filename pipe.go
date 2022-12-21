package cellnet

import (
	"sync"
)

// 不限制大小，添加不发生阻塞，接收阻塞等待
type Pipe struct {
	list      []interface{}
	listGuard sync.Mutex
	listCond  *sync.Cond
	ExistFlag int8
}

const DefaultExistFlag  = 1

// 添加时不会发送阻塞
func (self *Pipe) Add(msg interface{}) {
	self.listGuard.Lock()
	self.list = append(self.list, msg)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *Pipe) Count() int {
	self.listGuard.Lock()
	defer self.listGuard.Unlock()
	return len(self.list)
}

func (self *Pipe) Reset() {
	self.listGuard.Lock()
	self.list = self.list[0:0]
	self.listGuard.Unlock()
}

// 如果没有数据，发生阻塞
func (self *Pipe) Pick(retList *[]interface{}) (exit bool) {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		if self.ExistFlag == DefaultExistFlag {
			self.listGuard.Unlock()
			return true
		}
		self.listCond.Wait()
	}

	// 复制出队列
	*retList = append(*retList, self.list)

	self.list = self.list[0:0]
	self.listGuard.Unlock()

	return
}

func NewPipe() *Pipe {
	self := &Pipe{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}


func (self *Pipe) Close() {
	self.listGuard.Lock()
	self.ExistFlag = DefaultExistFlag
	self.listCond.Broadcast()
	self.listGuard.Unlock()
}
