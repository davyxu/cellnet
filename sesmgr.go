package cellnet

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
	Count() int
}
