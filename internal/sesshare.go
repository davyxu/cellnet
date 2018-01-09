package internal

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type tagData struct {
	name string
	data interface{}
}

type SessionShare struct {
	tag      []tagData
	tagGuard sync.RWMutex
	id       int64
	// 归属的通讯端
	PeerShare *CommunicatePeer
}

// 取会话归属的通讯端
func (self *SessionShare) Peer() cellnet.Peer {
	return self.PeerShare.Peer()
}

func (self *SessionShare) GetTag(name string) (interface{}, bool) {

	self.tagGuard.RLock()
	defer self.tagGuard.RUnlock()

	for _, t := range self.tag {
		if t.name == name {
			return t.data, true
		}
	}

	return nil, false
}

func (self *SessionShare) SetTag(name string, v interface{}) {

	self.tagGuard.Lock()
	defer self.tagGuard.Unlock()

	for i, t := range self.tag {
		if t.name == name {
			self.tag[i] = tagData{name, v}
			return
		}
	}

	self.tag = append(self.tag, tagData{name, v})
}

func (self *SessionShare) ID() int64 {
	return self.id
}

func (self *SessionShare) SetID(id int64) {
	self.id = id
}
