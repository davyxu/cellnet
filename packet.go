package cellnet

// 普通封包
type Packet struct {
	MsgID uint32 // 消息ID
	Data  []byte
}

func (self Packet) ContextID() int {
	return int(self.MsgID)
}
