package udppeer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
	"sync/atomic"
	"time"
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

	exitSignal chan bool // 通知开始退出线程

	recvTimeout time.Duration // 接收超时

	heartBeat int64 // 有收到封包时，为1，定时检查后设置为0

	endWaitor sync.WaitGroup

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

	self.exitSignal <- true

	self.endWaitor.Wait()
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
func (self *udpSession) Send(data interface{}) {

	raw := self.peer.CallOutboundProc(&cellnet.SendMsgEvent{self, data})
	if raw != nil {
		self.peer.CallInboundProc(&cellnet.SendMsgErrorEvent{self, raw.(error), data})

		self.Close()
	}
}

func (self *udpSession) HeartBeat() {

	atomic.StoreInt64(&self.heartBeat, 1)
}

func (self *udpSession) OnRecv(data []byte) error {

	self.HeartBeat()

	raw := self.peer.CallInboundProc(&cellnet.RecvDataEvent{self, data})
	if err, ok := raw.(error); ok && err != nil {
		self.peer.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})

		return err
	}

	return nil
}

// 启动会话的各种资源
func (self *udpSession) Start() {

	// 将会话添加到管理器
	self.Peer().(internal.SessionManager).Add(self)

	go func() {

		self.endWaitor.Add(1)

		for {

			select {
			case <-self.exitSignal:
				goto OnExit
			case <-time.After(self.recvTimeout):

				var targetValue int64
				currValue := atomic.SwapInt64(&self.heartBeat, targetValue)

				if currValue == 0 {
					self.peer.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})
					goto OnExit
				}

			}

		}

	OnExit:

		// 将会话从管理器移除
		self.Peer().(internal.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}

		self.endWaitor.Done()

	}()

}

func newUDPSession(addr *net.UDPAddr, conn *net.UDPConn, peer *internal.PeerShare, endNotify func()) *udpSession {
	return &udpSession{
		conn:        conn,
		remote:      addr,
		peer:        peer,
		recvTimeout: time.Second * 3,
		endNotify:   endNotify,
		exitSignal:  make(chan bool),
	}
}
