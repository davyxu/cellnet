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
	internal.SessionShare

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn

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

func (self *udpSession) Close() {

	self.exitSignal <- true

	self.endWaitor.Wait()
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

	raw := self.PeerShare.CallOutboundProc(&cellnet.SendMsgEvent{self, data})
	if raw != nil {
		self.PeerShare.CallInboundProc(&cellnet.SendMsgErrorEvent{self, raw.(error), data})

		self.Close()
	}
}

func (self *udpSession) HeartBeat() {

	atomic.StoreInt64(&self.heartBeat, 1)
}

func (self *udpSession) OnRecv(data []byte) error {

	self.HeartBeat()

	raw := self.PeerShare.CallInboundProc(&cellnet.RecvDataEvent{self, data})
	if err, ok := raw.(error); ok && err != nil {
		self.PeerShare.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})

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
					self.PeerShare.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})
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

func newUDPSession(addr *net.UDPAddr, conn *net.UDPConn, peerShare *internal.PeerShare, endNotify func()) *udpSession {
	self := &udpSession{
		conn:        conn,
		remote:      addr,
		recvTimeout: time.Second * 3,
		endNotify:   endNotify,
		exitSignal:  make(chan bool),
	}

	self.PeerShare = peerShare

	return self
}
