package cellnet

type Peer interface {

	// 开启
	Start(address string) Peer

	// 关闭
	Stop()

	// 名字
	SetName(string)
	Name() string

	// 路由模式, 收发封包中增加客户端标识号
	SetRelayMode(relay bool)

	// 事件
	EventManager

	// 连接管理
	SessionManager
}

type EventManager interface {

	// 注册事件回调
	RegisterCallback(id int, f func(interface{}))

	// 截获所有的事件
	Inject(func(interface{}) bool)
}

type SessionManager interface {

	// 获取一个连接
	Get(int64) Session

	// 遍历连接
	Iterate(func(Session) bool)

	// 连接数量
	Count() int
}
