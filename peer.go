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

type Peer interface {

	// 开启/关闭
	Start(address string) Peer

	Stop()

	// 名字
	SetName(string)
	Name() string

	// 地址
	SetAddress(string)
	Address() string

	// Tag
	SetTag(interface{})
	Tag() interface{}

	// 派发器
	EventDispatcher

	// 连接管理
	SessionManager

	//  HandlerList
	SetHandlerList(recv, send []EventHandler)

	HandlerList() (recv, send []EventHandler)

	Queue() EventQueue
}

// 连接器, 可由Peer转换
type Connector interface {

	// 连接后的Session
	DefaultSession() Session

	// 自动重连间隔, 0表示不重连, 默认不重连
	SetAutoReconnectSec(sec int)
}

// 会话管理器
type SessionManager interface {

	// 获取一个连接
	GetSession(int64) Session

	// 遍历连接
	VisitSession(func(Session) bool)

	// 连接数量
	SessionCount() int

	// 关闭所有连接
	CloseAllSession()
}
