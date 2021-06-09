package udp

import (
	cellpeer "github.com/davyxu/cellnet/peer"
	xframe "github.com/davyxu/x/frame"
	"net"
	"sync"
	"time"
)

// Socket会话
type Session struct {
	xframe.PropertySet
	cellpeer.SessionIdentify

	peer   *Peer
	parent interface{}

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
func (self *Session) Send(msg interface{}) {
	if self.peer.Recv == nil {
		panic("no transmitter")
	}

	// 在用户线程编码, 保证字段不会在其他线程被序列化读取
	ev := cellpeer.PackEvent(msg, &self.PropertySet)
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
