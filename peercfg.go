package cellnet

type PeerConfigInterface interface {
	// Peer的名称
	Name() string

	// Peer的类型，例如tcp.Connector
	TypeName() string

	// 侦听或者连接的地址:端口
	Address() string

	// 队列
	EventQueue() EventQueue
}

type PeerConfig struct {
	PeerType       string
	PeerName       string
	PeerAddress    string
	EventProcessor string

	Queue EventQueue

	UserInboundProc  EventProc
	UserOutboundProc EventProc
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
