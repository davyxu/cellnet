package cellnet

import (
	"sync"
	"sync/atomic"
)

// 会话访问
type SessionAccessor interface {

	// 获取一个连接
	GetSession(int64) Session

	// 遍历连接
	VisitSession(func(Session) bool)

	// 连接数量
	SessionCount() int

	// 关闭所有连接
	CloseAllSession()
}

// 完整功能的会话管理
type SessionManager interface {
	SessionAccessor

	Add(Session)
	Remove(Session)
}

type SessionManagerImplement struct {
	sesMap map[int64]Session

	sesIDAcc    int64
	sesMapGuard sync.RWMutex
}

const totalTryCount = 100

func (self *SessionManagerImplement) Add(ses Session) {

	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	var tryCount int = totalTryCount

	var id int64

	// id翻越处理
	for tryCount > 0 {

		id = atomic.AddInt64(&self.sesIDAcc, 1)

		if _, ok := self.sesMap[id]; !ok {
			break
		}

		tryCount--
	}

	if tryCount == 0 {
		log.Warnln("sessionID override!", id)
	}

	ses.(interface {
		SetID(int64)
	}).SetID(id)

	self.sesMap[id] = ses

}

func (self *SessionManagerImplement) Remove(ses Session) {
	self.sesMapGuard.Lock()
	delete(self.sesMap, ses.ID())
	self.sesMapGuard.Unlock()
}

// 获得一个连接
func (self *SessionManagerImplement) GetSession(id int64) Session {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	v, ok := self.sesMap[id]
	if ok {
		return v
	}

	return nil
}

func (self *SessionManagerImplement) VisitSession(callback func(Session) bool) {
	self.sesMapGuard.RLock()
	defer self.sesMapGuard.RUnlock()

	for _, ses := range self.sesMap {
		if !callback(ses) {
			break
		}
	}

}

func (self *SessionManagerImplement) CloseAllSession() {

	self.VisitSession(func(ses Session) bool {

		ses.Close()

		return true
	})
}

func (self *SessionManagerImplement) SessionCount() int {
	self.sesMapGuard.Lock()
	defer self.sesMapGuard.Unlock()

	return len(self.sesMap)
}

func NewSessionManager() SessionManager {
	return &SessionManagerImplement{
		sesMap: make(map[int64]Session),
	}
}
