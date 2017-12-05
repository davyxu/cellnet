package cellnet

type EventParam interface{}
type EventResult interface{}

// 事件函数的定义
type EventFunc func(EventParam) EventResult

// 发起和接受连接的通讯端
type Peer interface {

	// 开启端，传入地址
	Start() Peer

	// 停止通讯端
	Stop()

	SetEventFunc(EventFunc)

	EventFunc() EventFunc

	// 获取队列
	EventQueue() EventQueue

	// 通讯端名称
	Name() string

	Address() string

	TypeName() string

	SessionAccessor
}

type PeerConfig struct {
	PeerType    string
	PeerName    string
	PeerAddress string
	Queue       EventQueue
	Event       EventFunc
}

// 获取通讯端的名称
func (self *PeerConfig) Name() string {
	return self.PeerName
}

func (self *PeerConfig) TypeName() string {
	return self.PeerType
}

func (self *PeerConfig) Address() string {
	return self.PeerAddress
}

// 获取队列
func (self *PeerConfig) EventQueue() EventQueue {
	return self.Queue
}

func (self *PeerConfig) EventFunc() EventFunc {
	return self.Event
}

func (self *PeerConfig) SetEventFunc(eventFunc EventFunc) {
	self.Event = eventFunc
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

	f := creatorByTypeName[config.PeerType]
	if f == nil {
		panic("Peer name not found: " + config.PeerType)
	}

	return f(config)
}
