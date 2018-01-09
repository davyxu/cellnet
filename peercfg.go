package cellnet

type CommunicatePeerConfig struct {
	PeerType       string
	PeerName       string
	PeerAddress    string
	EventProcessor string

	Queue EventQueue

	UserInboundProc  EventProc
	UserOutboundProc EventProc
}

type PeerInfoSetter interface {
	SetAddress(addr string)
	SetName(name string)
	SetQueue(q EventQueue)
	SetEventFunc(processor string, inboundEvent, outboundEvent EventProc)
}

func CreatePeer(config CommunicatePeerConfig) Peer {

	p := NewPeer(config.PeerType)

	setter := p.(PeerInfoSetter)
	setter.SetName(config.PeerName)
	setter.SetAddress(config.PeerAddress)
	setter.SetQueue(config.Queue)
	setter.SetEventFunc(config.EventProcessor, config.UserInboundProc, config.UserOutboundProc)

	return p
}
