package peer

import (
	"net"
	"time"
)

type CoreTCPSocketOption struct {
	readBufferSize  int
	writeBufferSize int
	noDelay         bool
	maxPacketSize   int

	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (self *CoreTCPSocketOption) SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool) {
	self.readBufferSize = readBufferSize
	self.writeBufferSize = writeBufferSize
	self.noDelay = noDelay
}

func (self *CoreTCPSocketOption) SetSocketDeadline(read, write time.Duration) {

	self.readTimeout = read
	self.writeTimeout = write
}

func (self *CoreTCPSocketOption) SetMaxPacketSize(maxSize int) {
	self.maxPacketSize = maxSize
}

func (self *CoreTCPSocketOption) MaxPacketSize() int {

	return self.maxPacketSize
}

func (self *CoreTCPSocketOption) ApplySocketOption(conn net.Conn) {

	if cc, ok := conn.(*net.TCPConn); ok {

		if self.readBufferSize >= 0 {
			cc.SetReadBuffer(self.readBufferSize)
		}

		if self.writeBufferSize >= 0 {
			cc.SetWriteBuffer(self.writeBufferSize)
		}

		cc.SetNoDelay(self.noDelay)
	}

}

func (self *CoreTCPSocketOption) ApplySocketReadTimeout(conn net.Conn, callback func()) {

	if self.readTimeout > 0 {

		// issue: http://blog.sina.com.cn/s/blog_9be3b8f10101lhiq.html
		conn.SetReadDeadline(time.Now().Add(self.readTimeout))
		callback()
		conn.SetReadDeadline(time.Time{})

	} else {
		callback()
	}
}

func (self *CoreTCPSocketOption) ApplySocketWriteTimeout(conn net.Conn, callback func()) {

	if self.writeTimeout > 0 {

		conn.SetWriteDeadline(time.Now().Add(self.writeTimeout))
		callback()
		conn.SetWriteDeadline(time.Time{})

	} else {
		callback()
	}
}

func (self *CoreTCPSocketOption) Init() {
	self.readBufferSize = -1
	self.writeBufferSize = -1
}
