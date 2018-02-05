package peer

import "github.com/davyxu/cellnet"

type CommunicateConfig struct {
	PeerType       string
	PeerName       string
	PeerAddress    string
	EventProcessor string

	UserQueue cellnet.EventQueue

	UserInboundProc  cellnet.EventProc
	UserOutboundProc cellnet.EventProc
}

// 获取通讯端的名称
func (self *CommunicateConfig) Name() string {
	return self.PeerName
}

// 获取队列
func (self *CommunicateConfig) Queue() cellnet.EventQueue {
	return self.UserQueue
}
func (self *CommunicateConfig) Address() string {
	return self.PeerAddress
}

func (self *CommunicateConfig) SetAddress(addr string) {
	self.PeerAddress = addr
}

func (self *CommunicateConfig) NameOrAddress() string {
	if self.PeerName != "" {
		return self.PeerName
	}

	return self.PeerAddress
}

func (self *CommunicateConfig) InitConfig(config CommunicateConfig) {
	*self = config
}

func CreatePeer(config CommunicateConfig) cellnet.Peer {

	p := NewPeer(config.PeerType)

	type peerInitor interface {
		InitConfig(config CommunicateConfig)
		SetEventFunc(processor string, inboundEvent, outboundEvent cellnet.EventProc)
	}

	initor := p.(peerInitor)

	initor.InitConfig(config)
	initor.SetEventFunc(config.EventProcessor, config.UserInboundProc, config.UserOutboundProc)

	return p
}
