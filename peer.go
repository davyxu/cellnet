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

	IsConnector() bool

	IsAcceptor() bool

	PeerConfigInterface

	SessionAccessor
}

type PeerCreateFunc func() Peer

var creatorByTypeName = map[string]PeerCreateFunc{}

func RegisterPeerCreator(typeName string, f PeerCreateFunc) {

	if _, ok := creatorByTypeName[typeName]; ok {
		panic("Duplicate peer type")
	}

	creatorByTypeName[typeName] = f
}

type EventProcessor func(EventFunc, EventFunc) (EventFunc, EventFunc)

var evtprocByName = map[string]EventProcessor{}

func RegisterEventProcessor(name string, f EventProcessor) {

	if _, ok := creatorByTypeName[name]; ok {
		panic("Duplicate peer type")
	}

	evtprocByName[name] = f
}

func FetchEventProcessor(name string, inbound, outbound EventFunc) (EventFunc, EventFunc) {

	f := evtprocByName[name]
	if f == nil {
		panic("Event processor not found: " + name)
	}

	return f(inbound, outbound)
}

func CreatePeer(config PeerConfig) Peer {

	peerCreator := creatorByTypeName[config.PeerType]
	if peerCreator == nil {
		panic("Peer name not found: " + config.PeerType)
	}

	p := peerCreator()

	setter := p.(interface {
		SetConfig(PeerConfig)
		// 设置事件回调（入站，出站）
		SetEventFunc(inboundEvent, outboundEvent EventFunc)
	})

	setter.SetConfig(config)

	inboundEvent, outboundEvent := FetchEventProcessor(config.EventProcessor, config.InboundEvent, config.OutboundEvent)
	setter.SetEventFunc(inboundEvent, outboundEvent)

	return p
}
