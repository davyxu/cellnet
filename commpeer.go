package cellnet

import "sync"

// 通信端接口，例如tcp/udp的Peer
type CommunicatePeer interface {
	SetEventFunc(processor string, inboundEvent, outboundEvent EventProc)

	IsRunning() bool
}

// 通信通讯端共享的数据
type CoreCommunicatePeer struct {
	SessionManager
	CoreTagger
	CoreDuplexEventProc

	// 单独保存的保存Peer接口
	peerInterface Peer
	// 运行状态
	running      bool
	runningGuard sync.RWMutex

	// 停止过程同步
	stopping chan bool
}

func (self *CoreCommunicatePeer) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *CoreCommunicatePeer) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
}

func (self *CoreCommunicatePeer) Peer() Peer {
	return self.peerInterface
}

func (self *CoreCommunicatePeer) Init(p Peer) {
	self.SessionManager = NewSessionManager()
	self.peerInterface = p
}

func (self *CoreCommunicatePeer) WaitStopFinished() {
	// 如果正在停止时, 等待停止完成
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}

func (self *CoreCommunicatePeer) IsStopping() bool {
	return self.stopping != nil
}

func (self *CoreCommunicatePeer) StartStopping() {
	self.stopping = make(chan bool)
}

func (self *CoreCommunicatePeer) EndStopping() {
	select {
	case self.stopping <- true:

	default:
		self.stopping = nil
	}
}

type CommunicatePeerConfig struct {
	PeerType       string
	PeerName       string
	PeerAddress    string
	EventProcessor string

	Queue EventQueue

	UserInboundProc  EventProc
	UserOutboundProc EventProc
}

func CreatePeer(config CommunicatePeerConfig) Peer {

	p := NewPeer(config.PeerType)

	infoSetter := p.(PeerInfo)
	infoSetter.SetName(config.PeerName)
	infoSetter.SetAddress(config.PeerAddress)
	infoSetter.SetQueue(config.Queue)

	p.(CommunicatePeer).SetEventFunc(config.EventProcessor, config.UserInboundProc, config.UserOutboundProc)

	return p
}
