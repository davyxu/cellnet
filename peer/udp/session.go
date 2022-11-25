package udp

import (
	cellpeer "github.com/davyxu/cellnet/peer"
	"github.com/davyxu/x/container"
	"net"
	"sync"
	"time"
)

// Socket会话
type Session struct {
	xcontainer.Mapper
	cellpeer.SessionIdentify

	peer   *Peer
	parent any

	// Socket原始连接
	remote      *net.UDPAddr
	conn        *net.UDPConn
	connGuard   sync.RWMutex
	timeOutTick time.Time
	key         *connTrackKey
}

func (self *Session) IsAlive() bool {
	return time.Now().Before(self.timeOutTick)
}

func (self *Session) Read(data []byte) {

	if self.peer.Recv == nil {
		panic("no transmitter")
	}

	ev, err := self.peer.Recv(self, data)

	if ev != nil && err == nil {
		self.peer.ProcEvent(ev)
	}
}

func (self *Session) Write(data []byte) {

	// Connector中的Session
	if self.remote == nil {

		self.conn.Write(data)

		// Acceptor中的Session
	} else {
		self.conn.WriteToUDP(data, self.remote)
	}
}

// 发送封包
func (self *Session) Send(msg any) {
	if self.peer.Recv == nil {
		panic("no transmitter")
	}

	// 在用户线程编码, 保证字段不会在其他线程被序列化读取
	ev := cellpeer.PackEvent(msg, &self.Mapper)
	if ev == nil {
		return
	}
	ev.Ses = self

	if self.peer.OnOutbound != nil {
		ev = self.peer.OnOutbound(ev)
	}

	self.peer.Send(self, ev)
}

func (self *Session) Close() {

}
