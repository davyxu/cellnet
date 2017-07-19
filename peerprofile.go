package cellnet

import (
	"sync"
)

type PeerProfile interface {
	// 名字
	SetName(string)
	Name() string

	// 地址
	SetAddress(string)
	Address() string

	// Tag
	SetTag(interface{})
	Tag() interface{}
}

// Peer间的共享数据
type PeerProfileImplement struct {
	// 基本信息
	name    string
	address string
	tag     interface{}

	// 运行状态
	running      bool
	runningGuard sync.RWMutex
}

func (self *PeerProfileImplement) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *PeerProfileImplement) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
}

func (self *PeerProfileImplement) NameOrAddress() string {
	if self.name != "" {
		return self.name
	}

	return self.address
}

func (self *PeerProfileImplement) Tag() interface{} {
	return self.tag
}

func (self *PeerProfileImplement) SetTag(tag interface{}) {
	self.tag = tag
}

func (self *PeerProfileImplement) Address() string {
	return self.address
}

func (self *PeerProfileImplement) SetAddress(address string) {
	self.address = address
}

func (self *PeerProfileImplement) SetName(name string) {
	self.name = name
}

func (self *PeerProfileImplement) Name() string {
	return self.name
}

func NewPeerProfile() *PeerProfileImplement {

	return &PeerProfileImplement{}
}
