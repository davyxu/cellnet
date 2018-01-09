package internal

import (
	"github.com/davyxu/cellnet"
	"sync"
)

// 通信通讯端共享的数据
type CommunicatePeer struct {
	SessionManager

	// 单独保存的保存cellnet.Peer接口
	peerInterface cellnet.Peer
	// 运行状态
	running      bool
	runningGuard sync.RWMutex

	// 停止过程同步
	stopping chan bool

	InboundProc  cellnet.EventProc
	OutboundProc cellnet.EventProc
}

func (self *CommunicatePeer) SetEventFunc(processor string, inboundEvent, outboundEvent cellnet.EventProc) {
	self.InboundProc, self.OutboundProc = cellnet.MakeEventProcessor(processor, inboundEvent, outboundEvent)
}

func (self *CommunicatePeer) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *CommunicatePeer) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
}

// socket包内部派发事件
func (self *CommunicatePeer) CallInboundProc(ev interface{}) interface{} {

	if self.InboundProc == nil {
		return nil
	}

	//log.Debugf("<Inbound> %T|%+v", ev, ev)

	return self.InboundProc(ev)
}

// socket包内部派发事件
func (self *CommunicatePeer) CallOutboundProc(ev interface{}) interface{} {

	if self.OutboundProc == nil {
		return nil
	}

	//log.Debugf("<Outbound> %T|%+v", ev, ev)

	return self.OutboundProc(ev)
}

func (self *CommunicatePeer) Peer() cellnet.Peer {
	return self.peerInterface
}

func (self *CommunicatePeer) Init(p cellnet.Peer) {
	self.SessionManager = NewSessionManager()
	self.peerInterface = p
}

func (self *CommunicatePeer) WaitStopFinished() {
	// 如果正在停止时, 等待停止完成
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}

func (self *CommunicatePeer) IsStopping() bool {
	return self.stopping != nil
}

func (self *CommunicatePeer) StartStopping() {
	self.stopping = make(chan bool)
}

func (self *CommunicatePeer) EndStopping() {
	select {
	case self.stopping <- true:

	default:
		self.stopping = nil
	}
}
