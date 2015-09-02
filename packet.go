package cellnet

// 普通封包
type Packet struct {
	MsgID uint32 // 消息ID
	Data  []byte
}

// 封包流
type PacketStream interface {
	Read() (*Packet, error)
	Write(pkt *Packet) error
	Close() error
}
