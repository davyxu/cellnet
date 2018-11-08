package cellnet

import "time"

// UDP连接器
type UDPConnector interface {
	GenericPeer

	// 默认会话
	Session() Session
}

// UDP接受器
type UDPAcceptor interface {

	// 底层使用TTL做session生命期管理，超时时间越短，内存占用越低
	SetSessionTTL(dur time.Duration)
}
