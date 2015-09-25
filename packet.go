package cellnet

// 普通封包
type Packet struct {
	MsgID    uint32 // 消息ID
	Data     []byte
	ClientID int64 // 路由时, 客户端标识号
}

func (self Packet) ContextID() int {
	return int(self.MsgID)
}
