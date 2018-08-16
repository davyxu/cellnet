package cellnet

// 端, 可通过接口查询获得更多接口支持,如PeerProperty,ContextSet, SessionAccessor
type Peer interface {
	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	TypeName() string
}

// Peer基础属性
type PeerProperty interface {
	Name() string

	Address() string

	Queue() EventQueue

	// 设置名称（可选）
	SetName(v string)

	// 设置Peer地址
	SetAddress(v string)

	// 设置Peer挂接队列（可选）
	SetQueue(v EventQueue)
}

// 设置和获取自定义属性
type ContextSet interface {
	SetContext(key interface{}, v interface{})

	GetContext(key interface{}) (interface{}, bool)

	FetchContext(key, valuePtr interface{}) bool
}

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
