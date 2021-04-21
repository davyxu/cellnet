package cellpeer

import (
	"net"
	"time"
)

type SocketOption struct {
	ReadBufferSize  int
	WriteBufferSize int
	NoDelay         bool
	MaxPacketSize   int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (self *SocketOption) ApplySocketOption(conn net.Conn) {

	if cc, ok := conn.(*net.TCPConn); ok {

		if self.ReadBufferSize >= 0 {
			cc.SetReadBuffer(self.ReadBufferSize)
		}

		if self.WriteBufferSize >= 0 {
			cc.SetWriteBuffer(self.WriteBufferSize)
		}

		cc.SetNoDelay(self.NoDelay)
	}

}

func (self *SocketOption) BeginApplyReadTimeout(conn net.Conn) bool {
	// issue: http://blog.sina.com.cn/s/blog_9be3b8f10101lhiq.html
	if self.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(self.ReadTimeout))
		return true
	}

	return false
}

func (self *SocketOption) EndApplyTimeout(conn net.Conn) {
	conn.SetWriteDeadline(time.Time{})
}

func (self *SocketOption) BeginApplyWriteTimeout(conn net.Conn) bool {
	if self.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(self.WriteTimeout))
		return true
	}

	return false
}

func (self *SocketOption) Init() {
	self.ReadBufferSize = -1
	self.WriteBufferSize = -1
}
