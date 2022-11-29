package cellpeer

import (
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

type SessionIdentify struct {
	id int64
}

func (self *SessionIdentify) Id() int64 {
	return self.id
}

func (self *SessionIdentify) SetId(id int64) {
	self.id = id
}

type SessionID64Fetcher interface {
	Id() int64
}

type SessionID64Setter interface {
	SetId(id int64)
}

type SessionManager struct {
	sesById      map[int64]cellnet.Session
	sesByIdGuard sync.RWMutex

	sesIdGen int64 // 记录已经生成的会话ID流水号
}

// 设置id起始数值
func (self *SessionManager) SetSessionIdBase(base int64) {
	atomic.StoreInt64(&self.sesIdGen, base)
}

// 活跃的连接数量
func (self *SessionManager) SessionCount() int {
	self.sesByIdGuard.RLock()
	defer self.sesByIdGuard.RUnlock()
	return len(self.sesById)
}

// 将会话添加到管理中
func (self *SessionManager) AddSession(ses cellnet.Session) {

	id := atomic.AddInt64(&self.sesIdGen, 1)

	ses.(SessionID64Setter).SetId(id)

	self.sesByIdGuard.Lock()
	self.sesById[id] = ses
	self.sesByIdGuard.Unlock()
}

// 将会话移除管理
func (self *SessionManager) RemoveSession(ses cellnet.Session) {

	id := ses.(SessionID64Fetcher).Id()

	self.sesByIdGuard.Lock()
	delete(self.sesById, id)
	self.sesByIdGuard.Unlock()
}

// 获得一个会话
func (self *SessionManager) GetSession(id int64) cellnet.Session {
	self.sesByIdGuard.RLock()
	defer self.sesByIdGuard.RUnlock()
	return self.sesById[id]
}

// 遍历所有的会话
func (self *SessionManager) VisitSession(callback func(cellnet.Session) bool) {

	self.sesByIdGuard.RLock()
	defer self.sesByIdGuard.RUnlock()
	for _, ses := range self.sesById {
		if !callback(ses) {
			break
		}
	}
}

// 关闭所有会话
func (self *SessionManager) CloseAllSession() {

	self.sesByIdGuard.RLock()
	defer self.sesByIdGuard.RUnlock()
	for _, ses := range self.sesById {
		ses.(interface {
			Close()
		}).Close()
	}
	self.sesById = map[int64]cellnet.Session{}

}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sesById: map[int64]cellnet.Session{},
	}
}
