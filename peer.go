package cellnet

// 发起和接受连接的通讯端，某些Peer实现了SessionAccessor
type Peer interface {

	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	TypeName() string
}
