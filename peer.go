package cellnet

// 会话
type Session interface {

	// 发包
	Send(interface{})

	// 直接发送封包
	RawSend(*Event)

	// 断开
	Close()

	// 标示ID
	ID() int64

	// 归属端
	FromPeer() Peer

	// 将一个用户数据保存在session
	SetTag(tag interface{})

	// 取出与session关联的用户数据
	Tag() interface{}

	// 取原始连接net.Conn
	RawConn() interface{}
}

// 端, Connector或Acceptor
type Peer interface {

	// 开启/关闭
	Start(address string) Peer

	Stop()

	Queue() EventQueue

	// 基础信息
	PeerProfile

	// 定制处理链
	HandlerChainManager

	// 会话管理
	SessionAccessor
}
