package cellnet

// 事件函数的定义
type EventFunc func(interface{}) interface{}

// 发起和接受连接的通讯端
type Peer interface {

	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	// 获取队列
	Queue() EventQueue

	// 通讯端名称
	Name() string

	SessionAccessor
}

type PeerConfig struct {
	TypeName string
	Name     string
	Address  string
	Queue    EventQueue
	Event    EventFunc
}

type PeerCreateFunc func(PeerConfig) Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

func RegisterPeerCreator(typeName string, f PeerCreateFunc) {

	if _, ok := creatorByTypeName[typeName]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[typeName] = f
}

func NewPeer(config PeerConfig) Peer {

	f := creatorByTypeName[config.TypeName]
	if f == nil {
		panic("Peer name not found: " + config.TypeName)
	}

	return f(config)
}
