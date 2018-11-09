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

// 基本的通用Peer
type GenericPeer interface {
	Peer
	PeerProperty
}

// 设置和获取自定义属性
type ContextSet interface {
	// 为对象设置一个自定义属性
	SetContext(key interface{}, v interface{})

	// 从对象上根据key获取一个自定义属性
	GetContext(key interface{}) (interface{}, bool)

	// 给定一个值指针, 自动根据值的类型GetContext后设置到值
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

// 检查Peer是否正常工作
type PeerReadyChecker interface {
	IsReady() bool
}

// 开启IO层异常捕获,在生产版本对外端口应该打开此设置
type PeerCaptureIOPanic interface {
	// 开启IO层异常捕获
	EnableCaptureIOPanic(v bool)

	// 获取当前异常捕获值
	CaptureIOPanic() bool
}
