package cellnet

type Session interface {
	Send(interface{})

	Close()

	ID() int64

	FromPeer() Peer
}

type SessionManager interface {

	// 获取一个连接
	Get(int64) Session

	// 广播
	Broardcast(interface{})

	// 遍历连接
	Iterate(func(Session) bool)

	// 连接数量
	Count() int
}
