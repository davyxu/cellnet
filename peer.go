package cellnet

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

func CreatePeer(config PeerConfig) Peer {

	peerCreator := creatorByTypeName[config.PeerType]
	if peerCreator == nil {
		panic("Peer name not found: " + config.PeerType)
	}

	p := peerCreator()

	setter := p.(interface {
		SetConfig(PeerConfig)
		// 设置事件回调（入站，出站）
		SetEventFunc(inboundEvent, outboundEvent EventProc)
	})

	setter.SetConfig(config)

	inboundEvent, outboundEvent := FetchEventProcessor(config.EventProcessor, config.UserInboundProc, config.UserOutboundProc)
	setter.SetEventFunc(inboundEvent, outboundEvent)

	return p
}
