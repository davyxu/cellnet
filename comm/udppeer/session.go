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

type UPDSession interface {

	// 保活
	KeepAlive()

	// 设置接收超时
	SetRecvTimeout(duration time.Duration)

	// 标记会话经过验证
	MarkVerified()

	// 直接关闭，不通知
	RawClose()
}

// Socket会话
type udpSession struct {
	internal.SessionShare

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn

	exitSignal chan bool // 通知开始退出线程

	recvTimeout time.Duration // 接收超时

	recvBuffer []byte

	heartBeat int64 // 有收到封包时，为1，定时检查后设置为0

	verify int64

	closeState int64

	endWaitor sync.WaitGroup

	endNotify func()
}

// 取原始连接
func (self *udpSession) Raw() interface{} {
	return nil
}

func (self *udpSession) SetRecvTimeout(duration time.Duration) {
	self.recvTimeout = duration
}

type sendNotifier struct {
	data     interface{}
	callback func()
}

func (self *udpSession) Close() {

	self.Send(sendNotifier{data: &comm.SessionCloseNotify{}, callback: func() {
		self.RawClose()
	},
	})
}

func (self *udpSession) RawClose() {
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

	// 异步发送
	go func() {

		notifier, ok := data.(sendNotifier)

		if ok {
			data = notifier.data
		}

		raw := self.PeerShare.CallOutboundProc(&cellnet.SendMsgEvent{self, data})

		if raw != nil {

			if err, ok := raw.(error); ok {
				self.PeerShare.CallInboundProc(&cellnet.SendMsgErrorEvent{self, err, data})

				self.Close()
			}
		}

		if ok {
			notifier.callback()
		}

	}()
}

func (self *udpSession) KeepAlive() {

	atomic.StoreInt64(&self.heartBeat, 1)
}

func (self *udpSession) MarkVerified() {

	atomic.StoreInt64(&self.verify, 1)
}

func (self *udpSession) OnRecv(data []byte) {

	// 将数据拷贝到session的缓冲区
	self.recvBuffer = self.recvBuffer[0:len(data)]
	copy(self.recvBuffer, data)

	self.KeepAlive()
}

func (self *udpSession) ProcPacket() error {
	raw := self.PeerShare.CallInboundProc(&cellnet.RecvDataEvent{self, self.recvBuffer})
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

		var notifyEvent = true

		for {

			select {
			case <-self.exitSignal:
				goto OnExit
			case <-time.After(self.recvTimeout):

				verifyTag := atomic.LoadInt64(&self.verify)

				// 超时未验证
				if verifyTag == 0 {
					notifyEvent = false
					goto OnExit
				}

				var targetValue int64
				currValue := atomic.SwapInt64(&self.heartBeat, targetValue)

				if currValue == 0 {

					goto OnExit
				}

			}

		}

	OnExit:

		if notifyEvent {
			self.PeerShare.CallInboundProc(&cellnet.RecvMsgEvent{self, &comm.SessionClosed{}})
		}

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
		recvBuffer:  make([]byte, 0, MaxUDPRecvBuffer),
	}

	self.PeerShare = peerShare

	return self
}
