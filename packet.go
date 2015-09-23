package cellnet

import (
	"net"
)

// 普通封包
type Packet struct {
	MsgID uint32 // 消息ID
	Data  []byte
}

func (self Packet) ContextID() int {
	return int(self.MsgID)
}

// 封包流
type PacketStream interface {
	Read() (*Packet, error)
	Write(pkt *Packet) error
	Close() error
	Raw() net.Conn
}
