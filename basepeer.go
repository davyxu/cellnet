package cellnet

import (
	"sync"
)

type BasePeer interface {
	// 名字
	SetName(string)
	Name() string

	// 地址
	SetAddress(string)
	Address() string

	// Tag
	SetTag(interface{})
	Tag() interface{}

	//  HandlerList
	SetHandlerList(recv, send []EventHandler)

	HandlerList() (recv, send []EventHandler)
}

// Peer间的共享数据
type BasePeerImplement struct {
	// 基本信息
	name    string
	address string
	tag     interface{}

	// 接收, 发送处理器
	recvHandler  []EventHandler
	sendHandler  []EventHandler
	handlerGuard sync.RWMutex

	// 运行状态
	running      bool
	runningGuard sync.RWMutex
}

func (self *BasePeerImplement) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *BasePeerImplement) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
}

func (self *BasePeerImplement) NameOrAddress() string {
	if self.name != "" {
		return self.name
	}

	return self.address
}

func (self *BasePeerImplement) Tag() interface{} {
	return self.tag
}

func (self *BasePeerImplement) SetTag(tag interface{}) {
	self.tag = tag
}

func (self *BasePeerImplement) Address() string {
	return self.address
}

func (self *BasePeerImplement) SetAddress(address string) {
	self.address = address
}

func (self *BasePeerImplement) SetHandlerList(recv, send []EventHandler) {
	self.handlerGuard.Lock()
	self.recvHandler = recv
	self.sendHandler = send
	self.handlerGuard.Unlock()
}

func (self *BasePeerImplement) HandlerList() (recv, send []EventHandler) {
	self.handlerGuard.RLock()
	recv = self.recvHandler
	send = self.sendHandler
	self.handlerGuard.RUnlock()

	return
}

func (self *BasePeerImplement) SafeRecvHandler() (ret []EventHandler) {
	self.handlerGuard.RLock()
	ret = self.recvHandler
	self.handlerGuard.RUnlock()

	return
}

func (self *BasePeerImplement) SetName(name string) {
	self.name = name
}

func (self *BasePeerImplement) Name() string {
	return self.name
}

func NewBasePeer() *BasePeerImplement {

	return &BasePeerImplement{}
}
