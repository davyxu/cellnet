package udptransport

const (
	MTU       = 1472 // 最大传输单元
	packetLen = 2    // 包体大小字段
	MsgIDLen  = 2    // 消息ID字段

	HeaderSize = MsgIDLen + MsgIDLen // 整个UDP包头部分
)
