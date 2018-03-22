package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"net"
)

type DataReader interface {
	ReadData() []byte
}

type DataWriter interface {
	WriteData(data []byte)
}

// Socket会话
type udpSession struct {
	*peer.CoreProcBundle

	pInterface cellnet.Peer

	pkt []byte

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn
}

func (self *udpSession) ID() int64 {
	return 0
}

func (self *udpSession) Peer() cellnet.Peer {
	return self.pInterface
}

// 取原始连接
func (self *udpSession) Raw() interface{} {
	return self
}

func (self *udpSession) Recv(data []byte) {

	self.pkt = data

	msg, err := self.ReadMessage(self)

	if msg != nil && err == nil {
		self.PostEvent(&cellnet.RecvMsgEvent{self, msg})
	}
}

func (self *udpSession) ReadData() []byte {
	return self.pkt
}

func (self *udpSession) WriteData(data []byte) {

	if self.conn == nil {
		return
	}

	// Connector中的Session
	if self.remote == nil {

		self.conn.Write(data)

		// Acceptor中的Session
	} else {
		self.conn.WriteToUDP(data, self.remote)
	}
}

// 发送封包
func (self *udpSession) Send(msg interface{}) {

	self.SendMessage(&cellnet.SendMsgEvent{self, msg})
}

func (self *udpSession) Close() {

}

func newUDPSession(addr *net.UDPAddr, conn *net.UDPConn, p cellnet.Peer) *udpSession {
	self := &udpSession{
		conn:       conn,
		remote:     addr,
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}

	return self
}
