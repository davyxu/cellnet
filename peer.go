package cellnet

// 端, 可通过接口查询获得更多接口支持,如PeerProperty,PropertySet, SessionAccessor
type Peer interface {
	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	TypeName() string
}

// Peer属性,基于PropertySet
type PeerProperty interface {
	Name() (ret string)
	Queue() (ret EventQueue)
	Address() (ret string)
}

// 设置和获取预制属性,自定义属性
type PropertySet interface {
	GetProperty(key, valuePtr interface{}) bool

	SetProperty(key interface{}, v interface{})

	RawGetProperty(key interface{}) (interface{}, bool)
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
