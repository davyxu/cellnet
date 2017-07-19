package cellnet

import (
	"bytes"
	"sync"
)

type HandlerChainManager interface {

	// 添加一条接收处理链
	AddChainRecv(recv *HandlerChain) int64

	// 移除接收处理链, 根据添加时的id
	RemoveChainRecv(id int64)

	// 获取当前的处理链(乱序)
	ChainListRecv() HandlerChainList

	// 设置发送处理链
	SetChainSend(chain *HandlerChain)

	// 获取当前发送处理链
	ChainSend() *HandlerChain
}

// Peer间的共享数据
type HandlerChainManagerImplement struct {
	recvChainByID      map[int64]*HandlerChain
	recvChainGuard     sync.Mutex
	chainIDAcc         int64
	recvChainListDirty bool
	recvChainList      HandlerChainList

	sendChain      *HandlerChain
	sendChainGuard sync.RWMutex
}

func (self *HandlerChainManagerImplement) AddChainRecv(recv *HandlerChain) (id int64) {

	self.recvChainGuard.Lock()

	self.chainIDAcc++
	id = self.chainIDAcc
	self.recvChainByID[id] = recv
	self.recvChainListDirty = true

	self.recvChainGuard.Unlock()

	return
}

func (self *HandlerChainManagerImplement) RemoveChainRecv(id int64) {

	self.recvChainGuard.Lock()

	delete(self.recvChainByID, id)
	self.recvChainListDirty = true

	self.recvChainGuard.Unlock()
}

func (self *HandlerChainManagerImplement) SetChainSend(chain *HandlerChain) {

	self.sendChainGuard.Lock()
	self.sendChain = chain
	self.sendChainGuard.Unlock()
}

func (self *HandlerChainManagerImplement) ChainSend() *HandlerChain {
	self.sendChainGuard.Lock()
	defer self.sendChainGuard.Unlock()
	return self.sendChain
}

func (self *HandlerChainManagerImplement) ChainListRecv() HandlerChainList {
	self.recvChainGuard.Lock()
	defer self.recvChainGuard.Unlock()

	if self.recvChainListDirty {

		self.recvChainList = self.recvChainList[0:0]
		for _, chain := range self.recvChainByID {
			self.recvChainList = append(self.recvChainList, chain)
		}

		self.recvChainListDirty = false
	}

	return self.recvChainList
}

func (self *HandlerChainManagerImplement) ChainString() string {

	var buff bytes.Buffer

	buff.WriteString("ChainRecv:\n")
	buff.WriteString(self.ChainListRecv().String())

	buff.WriteString("ChainSend:\n")
	buff.WriteString(self.ChainSend().String())

	return buff.String()
}

func NewHandlerChainManager() *HandlerChainManagerImplement {

	return &HandlerChainManagerImplement{
		recvChainByID: make(map[int64]*HandlerChain),
	}
}
