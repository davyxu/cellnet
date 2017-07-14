package cellnet

// 会话
type Session interface {

	// 发包
	Send(interface{})

	// 直接发送封包
	RawSend([]EventHandler, *Event)

	// 投递封包
	Post(interface{})

	// 直接投递封包
	RawPost([]EventHandler, *Event)

	// 断开
	Close()

	// 标示ID
	ID() int64

	// 归属端
	FromPeer() Peer
}

// 端, Connector或Acceptor
type Peer interface {

	// 开启/关闭
	Start(address string) Peer

	Stop()

	// 扩展用的功能
	BasePeer

	// 派发器
	EventDispatcher

	// 连接管理
	SessionAccessor

	Queue() EventQueue
}
