package peer

import (
	"sync"
)

// 消息队列，用于避免固定大小的channel实现的队列发生阻塞情况
type MsgQueue struct {
	list      []interface{}
	listGuard sync.Mutex
	listCond  *sync.Cond
}

func (self *MsgQueue) Add(msg interface{}) {
	self.listGuard.Lock()
	self.list = append(self.list, msg)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *MsgQueue) Reset() {
	self.list = self.list[0:0]
}

func (self *MsgQueue) Pick(retList *[]interface{}) (exit bool) {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	self.listGuard.Unlock()

	self.listGuard.Lock()

	// 复制出队列

	for _, ev := range self.list {

		if ev == nil {
			exit = true
			break
		} else {
			*retList = append(*retList, ev)
		}
	}

	self.Reset()
	self.listGuard.Unlock()

	return
}

func NewMsgQueue() *MsgQueue {
	self := &MsgQueue{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}
