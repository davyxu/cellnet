package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
)

// Socket会话
type udpSession struct {

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn

	tag interface{}

	// 归属的通讯端
	peer *internal.PeerShare

	id int64

	endNotify func()
}

// 取原始连接
func (self *udpSession) Raw() interface{} {
	return nil
}

func (self *udpSession) Tag() interface{} {
	return self.tag
}

func (self *udpSession) SetTag(v interface{}) {
	self.tag = v
}

func (self *udpSession) ID() int64 {
	return self.id
}

func (self *udpSession) SetID(id int64) {
	self.id = id
}

func (self *udpSession) Close() {

}

// 取会话归属的通讯端
func (self *udpSession) Peer() cellnet.Peer {
	return self.peer.Peer()
}

func (self *udpSession) WriteData(data []byte) error {

	if self.Peer().IsConnector() {

		_, err := self.conn.Write(data)
		return err

	} else {
		_, err := self.conn.WriteToUDP(data, self.remote)
		return err
	}
}

// 发送封包
func (self *udpSession) Send(msg interface{}) {

	raw := self.peer.FireEvent(cellnet.SendMsgEvent{self, msg})
	if raw != nil {
		self.peer.FireEvent(cellnet.SendMsgErrorEvent{self, raw.(error), msg})
	}
}

func (self *udpSession) OnRecv(data []byte) error {

	raw := self.peer.FireEvent(cellnet.RecvDataEvent{self, data})
	if err, ok := raw.(error); ok && err != nil {
		self.peer.FireEvent(cellnet.SessionClosedEvent{self, err})

		return err
	}

	return nil
}

// 启动会话的各种资源
func (self *udpSession) Start() {

	// 将会话添加到管理器
	self.Peer().(internal.SessionManager).Add(self)

	//// 将会话从管理器移除
	//self.Peer().(internal.SessionManager).Remove(self)
	//
	//if self.endNotify != nil {
	//	self.endNotify()
	//}
}

func newUDPSession(addr *net.UDPAddr, conn *net.UDPConn, peer *internal.PeerShare, endNotify func()) *udpSession {
	return &udpSession{
		conn:      conn,
		remote:    addr,
		peer:      peer,
		endNotify: endNotify,
	}
}
