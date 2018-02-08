package cellnet

// 长连接
type Session interface {

	// 获得原始的Socket连接
	Raw() interface{}

	// 获得Session归属的Peer
	Peer() Peer

	// 发送消息，消息需要以指针格式传入
	Send(msg interface{})

	// 断开
	Close()

	// 标示ID
	ID() int64
}
