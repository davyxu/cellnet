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

// 内部使用, 连接初始化时, 将配置应用到TCP连接上
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

// 每次读取前设置读取超时
func (self *SocketOption) BeginApplyReadTimeout(conn net.Conn) bool {
	// issue: http://blog.sina.com.cn/s/blog_9be3b8f10101lhiq.html
	if self.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(self.ReadTimeout))
		return true
	}

	return false
}

// 每次读取后, 恢复读取超时
func (self *SocketOption) EndApplyTimeout(conn net.Conn) {
	conn.SetWriteDeadline(time.Time{})
}

// 每次写入前设置写入超时
func (self *SocketOption) BeginApplyWriteTimeout(conn net.Conn) bool {
	if self.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(self.WriteTimeout))
		return true
	}

	return false
}

// 连接构造时, 初始化参数
func (self *SocketOption) Init() {
	self.ReadBufferSize = -1
	self.WriteBufferSize = -1
}
