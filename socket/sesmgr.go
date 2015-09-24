package socket

import (
	"github.com/davyxu/cellnet"
	"log"
	"sync"
	"sync/atomic"
)

type sessionMgr struct {
	sesMap map[uint32]cellnet.Session

	sesIDAcc    uint32
	sesMapGuard sync.RWMutex
}

const totalTryCount = 100

func (self *sessionMgr) Add(ses cellnet.Session) {

	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	var tryCount int = totalTryCount

	var id uint32

	// id翻越处理
	for tryCount > 0 {

		id = atomic.AddUint32(&self.sesIDAcc, 1)

		if _, ok := self.sesMap[id]; !ok {
			break
		}

		tryCount--
	}

	if tryCount == 0 {
		log.Println("WARNING: sessionID override!", id)
	}

	ltvses := ses.(*ltvSession)

	ltvses.id = id

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
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	for _, ses := range self.sesMap {
		ses.Send(data)
	}

}

// 获得一个连接
func (self *sessionMgr) Get(id uint32) cellnet.Session {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	v, ok := self.sesMap[id]
	if ok {
		return v
	}

	return nil
}

func (self *sessionMgr) Iterate(callback func(cellnet.Session) bool) {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

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
		sesMap: make(map[uint32]cellnet.Session),
	}
}
