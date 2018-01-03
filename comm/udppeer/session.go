package udppeer

import (
	"errors"
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
	RawClose(err error)
}

// Socket会话
type udpSession struct {
	internal.SessionShare

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn

	exitSignal chan error // 通知开始退出线程

	recvTimeout time.Duration // 接收超时

	recvBuffer []byte

	heartBeat int64 // 有收到封包时，为1，定时检查后设置为0

	verify int64

	closing bool

	closeNotify bool

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

type sendDeferCall struct {
	data     interface{}
	callback func()
}

func (self *udpSession) Close() {

	if self.closeNotify {
		return
	}

	self.closeNotify = true

	self.Send(sendDeferCall{data: &comm.SessionCloseNotify{}, callback: func() {
		self.RawClose(nil)
	},
	})
}

func (self *udpSession) RawClose(err error) {

	if self.closing {
		return
	}

	self.closing = true

	self.exitSignal <- err

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

		notifier, ok := data.(sendDeferCall)

		if ok {
			data = notifier.data
		}

		raw := self.PeerShare.CallOutboundProc(&cellnet.SendMsgEvent{self, data})

		if err, ok := raw.(error); ok && err != nil && !ok {

			self.RawClose(err)
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

	if self.closing {
		return
	}

	// TODO 改造为sync.Pool

	// 将数据拷贝到session的缓冲区
	self.recvBuffer = self.recvBuffer[0:len(data)]
	copy(self.recvBuffer, data)

	self.KeepAlive()
}

func (self *udpSession) ProcPacket() error {
	raw := self.PeerShare.CallInboundProc(&cellnet.RecvDataEvent{self, self.recvBuffer})

	if err, ok := raw.(error); ok && err != nil {

		self.RawClose(err)

		return err
	}

	return nil
}

var (
	ErrNotVerify   = errors.New("UDPSession not verify")
	ErrReadTimeout = errors.New("UDPSession timeout")
)

// 启动会话的各种资源
func (self *udpSession) Start() {

	// 将会话添加到管理器
	self.Peer().(internal.SessionManager).Add(self)

	go func() {

		self.endWaitor.Add(1)

		var notifyEvent = true

		var err error

		for {

			select {
			case err = <-self.exitSignal:
				// 正常退出
				goto OnExit
			case <-time.After(self.recvTimeout):

				verifyTag := atomic.LoadInt64(&self.verify)

				// 超时未验证
				if verifyTag == 0 {
					err = ErrNotVerify
					notifyEvent = false
					goto OnExit
				}

				var targetValue int64
				currValue := atomic.SwapInt64(&self.heartBeat, targetValue)

				// 心跳超时
				if currValue == 0 {
					err = ErrReadTimeout
					goto OnExit
				}

			}

		}

	OnExit:

		if notifyEvent {
			msg := &comm.SessionClosed{}

			if err != nil {
				msg.Error = err.Error()
			}

			self.PeerShare.CallInboundProc(&cellnet.RecvMsgEvent{self, msg})
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
		exitSignal:  make(chan error),
		recvBuffer:  make([]byte, 0, MaxUDPRecvBuffer),
	}

	self.PeerShare = peerShare

	return self
}
