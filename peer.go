package cellnet

type Peer interface {

	// 开启
	Start(address string) Peer

	// 关闭
	Stop()

	// 名字
	SetName(string)
	Name() string

	// 事件
	EventQueue

	// 连接管理
	SessionManager
}

type SessionManager interface {

	// 获取一个连接
	GetSession(int64) Session

	// 遍历连接
	IterateSession(func(Session) bool)

	// 连接数量
	SessionCount() int
}
