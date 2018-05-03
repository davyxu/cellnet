package peer

import (
	"net"
)

type CoreTCPSocketOption struct {
	readBufferSize  int
	writeBufferSize int
	noDelay         bool
}

func (self *CoreTCPSocketOption) SetSocketBuffer(readBufferSize, writeBufferSize int, noDelay bool) {
	self.readBufferSize = readBufferSize
	self.writeBufferSize = writeBufferSize
	self.noDelay = noDelay
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

func (self *CoreTCPSocketOption) Init() {
	self.readBufferSize = -1
	self.writeBufferSize = -1
}
