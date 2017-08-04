package socket

import (
	"net"
	"time"
)

type SocketOptions interface {
	// Session最大包大小, 超过这个数字, 接收视为错误, 断开连接
	SetMaxPacketSize(size int)

	MaxPacketSize() int

	// 设置socket选项, 如果不修改,请设置-1
	SetSocketOption(readBufferSize, writeBufferSize int, nodelay bool)

	// 设置socket超时间隔, 0表示不作用
	SetSocketDeadline(read, write time.Duration)
	SocketDeadline() (read, write time.Duration)
}

type socketOptions struct {
	// socket参数
	maxPacketSize    int
	connReadBuffer   int
	connWriteBuffer  int
	connNoDelay      bool
	connReadTimeout  time.Duration
	connWriteTimeout time.Duration
}

// socket配置
func (self *socketOptions) SetMaxPacketSize(size int) {
	self.maxPacketSize = size
}

func (self *socketOptions) MaxPacketSize() int {
	return self.maxPacketSize
}

func (self *socketOptions) SetSocketDeadline(read, write time.Duration) {
	self.connReadTimeout = read
	self.connWriteTimeout = write
}

func (self *socketOptions) SocketDeadline() (read, write time.Duration) {
	return self.connReadTimeout, self.connWriteTimeout
}

func (self *socketOptions) SetSocketOption(readBufferSize, writeBufferSize int, nodelay bool) {

	self.connReadBuffer = readBufferSize
	self.connWriteBuffer = writeBufferSize
	self.connNoDelay = nodelay
}

func (self *socketOptions) Apply(conn net.Conn) {

	if cc, ok := conn.(*net.TCPConn); ok {

		if self.connReadBuffer >= 0 {
			cc.SetReadBuffer(self.connReadBuffer)
		}

		if self.connWriteBuffer >= 0 {
			cc.SetWriteBuffer(self.connWriteBuffer)
		}

		cc.SetNoDelay(self.connNoDelay)
	}

}

func newSocketOptions() *socketOptions {
	return &socketOptions{
		connWriteBuffer: -1,
		connReadBuffer:  -1,
	}
}
