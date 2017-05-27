package cellnet

import "net"

type PacketStream interface {
	Read() (msgid uint32, data []byte, err error)

	Write(msgid uint32, data []byte) error

	// 将Write写入的数据提交
	Flush() error

	// 关闭连接
	Close() error

	Raw() net.Conn
}
