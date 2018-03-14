package udp

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
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

type DataReader interface {
	ReadData() []byte
}

type DataWriter interface {
	WriteData(data []byte)
}

// Socket会话
type udpSession struct {
	peer.CorePropertySet
	peer.CoreSessionIdentify
	*peer.CoreProcessorBundle

	pInterface cellnet.Peer

	recvChan chan []byte

	// Socket原始连接
	remote *net.UDPAddr
	conn   *net.UDPConn

	exitSignal chan error // 通知开始退出线程

	recvTimeout time.Duration // 接收超时

	heartBeat int64 // 有收到封包时，为1，定时检查后设置为0

	verify int64

	closing bool

	closeNotify bool

	endWaitor sync.WaitGroup

	endNotify func()
}

func (self *udpSession) Peer() cellnet.Peer {
	return self.pInterface
}

// 取原始连接
func (self *udpSession) Raw() interface{} {
	return self
}

func (self *udpSession) SetRecvTimeout(duration time.Duration) {
	self.recvTimeout = duration
}

func (self *udpSession) Recv(data []byte) {

	if self.closing {
		return
	}

	self.KeepAlive()

	self.recvChan <- data
}

func (self *udpSession) recvLoop() {

	for self.conn != nil {

		msg, err := self.ReadMessage(self)

		if err != nil {
			self.PostEvent(&cellnet.RecvMsgEvent{self, &cellnet.SessionClosed{
				Error: err.Error(),
			}})
			break
		}

		if _, ok := msg.(*cellnet.SessionCloseNotify); ok {

			self.RawClose(nil)
		} else {

			self.PostEvent(&cellnet.RecvMsgEvent{self, msg})
		}

	}
}

func (self *udpSession) ReadData() []byte {
	return <-self.recvChan
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

	if self.closeNotify {
		return
	}

	self.closeNotify = true

	go func() {

		self.Send(&cellnet.SessionCloseNotify{})
		self.RawClose(nil)

	}()

}

func (self *udpSession) RawClose(err error) {

	if self.closing {
		return
	}

	self.closing = true

	self.conn = nil

	self.exitSignal <- err

	self.endWaitor.Wait()
}

func (self *udpSession) KeepAlive() {

	atomic.StoreInt64(&self.heartBeat, 1)
}

func (self *udpSession) MarkVerified() {

	atomic.StoreInt64(&self.verify, 1)
}

var (
	ErrNotVerify   = errors.New("UDPSession not verify")
	ErrReadTimeout = errors.New("UDPSession timeout")
)

func (self *udpSession) tickLoop() {
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
		msg := &cellnet.SessionClosed{}

		if err != nil {
			msg.Error = err.Error()
		}

		self.PostEvent(&cellnet.RecvMsgEvent{self, msg})
	}

	// 将会话从管理器移除
	self.pInterface.(peer.SessionManager).Remove(self)

	if self.endNotify != nil {
		self.endNotify()
	}

	self.endWaitor.Done()
}

// 启动会话的各种资源
func (self *udpSession) Start() {

	self.PostEvent(&cellnet.RecvMsgEvent{self, &cellnet.SessionInit{}})

	// 将会话添加到管理器
	self.pInterface.(peer.SessionManager).Add(self)

	go self.tickLoop()

	go self.recvLoop()
}

const RecvBufferLen = 10

func newUDPSession(addr *net.UDPAddr, conn *net.UDPConn, p cellnet.Peer, endNotify func()) *udpSession {
	self := &udpSession{
		conn:        conn,
		remote:      addr,
		recvTimeout: time.Second * 3,
		endNotify:   endNotify,
		exitSignal:  make(chan error),
		recvChan:    make(chan []byte, RecvBufferLen),
		pInterface:  p,
		CoreProcessorBundle: p.(interface {
			GetBundle() *peer.CoreProcessorBundle
		}).GetBundle(),
	}

	return self
}
