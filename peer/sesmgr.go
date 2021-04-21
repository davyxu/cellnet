package cellpeer

import (
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

type SessionIdentify struct {
	id int64
}

func (self *SessionIdentify) ID() int64 {
	return self.id
}

func (self *SessionIdentify) SetID(id int64) {
	self.id = id
}

type identifyChecker interface {
	ID() int64
	SetID(id int64)
}

type SessionManager struct {
	sesByID      map[int64]cellnet.Session
	sesByIDGuard sync.RWMutex

	sesIDGen int64 // 记录已经生成的会话ID流水号
}

func (self *SessionManager) SetIDBase(base int64) {
	atomic.StoreInt64(&self.sesIDGen, base)
}

func (self *SessionManager) Count() int {
	self.sesByIDGuard.RLock()
	defer self.sesByIDGuard.RUnlock()
	return len(self.sesByID)
}

func (self *SessionManager) Add(ses cellnet.Session) {

	id := atomic.AddInt64(&self.sesIDGen, 1)

	ses.(identifyChecker).SetID(id)

	self.sesByIDGuard.Lock()
	self.sesByID[id] = ses
	self.sesByIDGuard.Unlock()
}

func (self *SessionManager) Remove(ses cellnet.Session) {

	id := ses.(identifyChecker).ID()

	self.sesByIDGuard.Lock()
	delete(self.sesByID, id)
	self.sesByIDGuard.Unlock()
}

// 获得一个连接
func (self *SessionManager) Get(id int64) cellnet.Session {
	self.sesByIDGuard.RLock()
	defer self.sesByIDGuard.RUnlock()
	return self.sesByID[id]
}

func (self *SessionManager) Visit(callback func(cellnet.Session) bool) {

	self.sesByIDGuard.RLock()
	defer self.sesByIDGuard.RUnlock()
	for _, ses := range self.sesByID {
		if !callback(ses) {
			break
		}
	}
}

func (self *SessionManager) CloseAll() {

	self.sesByIDGuard.RLock()
	defer self.sesByIDGuard.RUnlock()
	for _, ses := range self.sesByID {
		ses.Close()
	}
	self.sesByID = map[int64]cellnet.Session{}

}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sesByID: map[int64]cellnet.Session{},
	}
}
