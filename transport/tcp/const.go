package tcptransport

const (
	packetHeaderSize = 2 // 包体大小字段
	msgIDLen         = 2 // 消息ID字段
)

var (
	TestEnableRecvPanic bool
	TestEnableSendPanic bool
)
