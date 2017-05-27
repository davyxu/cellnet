package cellnet

import "time"

type Session interface {

	// 发包
	Send(interface{})

	// 直接发送封包
	RawSend([]EventHandler, *SessionEvent)

	// 投递封包
	Post(interface{})

	// 直接投递封包
	RawPost([]EventHandler, *SessionEvent)

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

	// Session最大包大小, 超过这个数字, 接收视为错误, 断开连接
	SetMaxPacketSize(size int)

	MaxPacketSize() int

	// 设置socket选项, 如果不修改,请设置-1
	SetSocketOption(readBufferSize, writeBufferSize int, nodelay bool)

	// 设置socket超时间隔, 0表示不作用
	SetSocketDeadline(read, write time.Duration)
	SocketDeadline() (read, write time.Duration)

	// 派发器
	EventDispatcher

	// 连接管理
	SessionManager

	//  HandlerList
	SetHandlerList(recv, send []EventHandler)

	HandlerList() (recv, send []EventHandler)

	Queue() EventQueue
}

type Connector interface {

	// 连接后的Session
	DefaultSession() Session

	// 自动重连间隔, 0表示不重连, 默认不重连
	SetAutoReconnectSec(sec int)
}

type SessionManager interface {

	// 获取一个连接
	GetSession(int64) Session

	// 遍历连接
	VisitSession(func(Session) bool)

	// 连接数量
	SessionCount() int
}
