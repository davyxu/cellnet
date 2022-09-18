package tcptransport

const (
	packetHeaderSize = 4 // 包体大小字段
	msgIDLen         = 4 // 消息ID字段
)

var (
	TestEnableRecvPanic bool
	TestEnableSendPanic bool
)
