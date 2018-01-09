package cellnet

// 发起和接受连接的通讯端，某些Peer实现了SessionAccessor
type Peer interface {

	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	// Peer的名称(自定义名称，会出现在msglog中)
	Name() string

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	TypeName() string

	// 侦听或者连接的地址:端口
	Address() string

	// 队列
	EventQueue() EventQueue
}

type PeerCreateFunc func() Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

func RegisterPeerCreator(f PeerCreateFunc) {

	// 临时实例化一个，获取类型
	dummyPeer := f()

	if _, ok := creatorByTypeName[dummyPeer.TypeName()]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[dummyPeer.TypeName()] = f
}

func NewPeer(peerType string) Peer {
	peerCreator := creatorByTypeName[peerType]
	if peerCreator == nil {
		panic("Peer name not found: " + peerType)
	}

	return peerCreator()
}
