package cellnet

type Session interface {

	// 发包
	Send(interface{})

	// 直接发送封包
	RawSend(EventHandler, *SessionEvent)

	// 断开
	Close()

	// 标示ID
	ID() int64

	// 归属端
	FromPeer() Peer
}
