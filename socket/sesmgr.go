package socket

import (
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

type sessionMgr struct {
	sesMap map[int64]cellnet.Session

	sesIDAcc    int64
	sesMapGuard sync.Mutex
}

func (self *sessionMgr) Add(ses cellnet.Session) {

	var id int64

	ltvses := ses.(*ltvSession)

	id = atomic.AddInt64(&self.sesIDAcc, 1)

	ltvses.SetID(id)

	self.sesMapGuard.Lock()
	self.sesMap[id] = ses
	self.sesMapGuard.Unlock()

}

func (self *sessionMgr) Remove(ses cellnet.Session) {
	self.sesMapGuard.Lock()
	delete(self.sesMap, ses.ID())
	self.sesMapGuard.Unlock()
}

// 广播到所有连接
func (self *sessionMgr) Broardcast(data interface{}) {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	for _, ses := range self.sesMap {
		ses.Send(data)
	}

}

// 获得一个连接
func (self *sessionMgr) Get(id int64) cellnet.Session {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	v, ok := self.sesMap[id]
	if ok {
		return v
	}

	return nil
}

func (self *sessionMgr) Iterate(callback func(cellnet.Session) bool) {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	for _, ses := range self.sesMap {
		if !callback(ses) {
			break
		}
	}

}

func (self *sessionMgr) Count() int {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	return len(self.sesMap)
}

func newSessionManager() *sessionMgr {
	return &sessionMgr{
		sesMap: make(map[int64]cellnet.Session),
	}
}
