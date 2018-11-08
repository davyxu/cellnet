package cellnet

import "time"

// TCP
type TCPSocketOption interface {
	// 收发缓冲大小，默认-1
	SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool)

	// 设置最大的封包大小
	SetMaxPacketSize(maxSize int)

	// 设置读写超时，默认0，不超时
	SetSocketDeadline(read, write time.Duration)
}

// TCP接受器，具备会话访问
type TCPAcceptor interface {
	GenericPeer

	// 访问会话
	SessionAccessor

	TCPSocketOption

	// 查看当前侦听端口，使用host:0 作为Address时，socket底层自动分配侦听端口
	Port() int
}

// TCP连接器
type TCPConnector interface {
	GenericPeer

	TCPSocketOption

	// 设置重连时间
	SetReconnectDuration(time.Duration)

	// 获取重连时间
	ReconnectDuration() time.Duration

	// 默认会话
	Session() Session

	// 设置会话管理器 实现peer.SessionManager接口
	SetSessionManager(raw interface{})

	// 查看当前连接使用的端口
	Port() int
}
