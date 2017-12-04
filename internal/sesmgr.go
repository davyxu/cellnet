package internal

import (
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

// 完整功能的会话管理
type SessionManager interface {
	cellnet.SessionAccessor

	Add(cellnet.Session)
	Remove(cellnet.Session)
}

type sesMgr struct {
	sesById sync.Map // 使用Id关联会话

	sesIDGen int64 // 记录已经生成的会话ID流水号

	count int64 // 记录当前在使用的会话数量
}

func (self *sesMgr) Add(ses cellnet.Session) {

	id := atomic.AddInt64(&self.sesIDGen, 1)

	self.count = atomic.AddInt64(&self.count, 1)

	ses.(interface {
		SetID(int64)
	}).SetID(id)

	self.sesById.Store(id, ses)
}

func (self *sesMgr) Remove(ses cellnet.Session) {

	self.sesById.Delete(ses.ID())

	self.count = atomic.AddInt64(&self.count, -1)
}

// 获得一个连接
func (self *sesMgr) GetSession(id int64) cellnet.Session {
	if v, ok := self.sesById.Load(id); ok {
		return v.(cellnet.Session)
	}

	return nil
}

func (self *sesMgr) VisitSession(callback func(cellnet.Session) bool) {

	self.sesById.Range(func(key, value interface{}) bool {

		return callback(value.(cellnet.Session))

	})
}

func (self *sesMgr) CloseAllSession() {

	self.VisitSession(func(ses cellnet.Session) bool {

		ses.Close()

		return true
	})
}

func (self *sesMgr) SessionCount() int {

	v := atomic.LoadInt64(&self.count)

	return int(v)
}

func NewSessionManager() SessionManager {
	return &sesMgr{}
}
