package socket

import (
	"github.com/davyxu/cellnet"
	"net"
	"sync"
	"time"
)

// Peer间的共享数据
type peerBase struct {
	cellnet.EventQueue

	// 基本信息
	name    string
	address string
	tag     interface{}

	// socket参数
	maxPacketSize    int
	connReadBuffer   int
	connWriteBuffer  int
	connNoDelay      bool
	connReadTimeout  time.Duration
	connWriteTimeout time.Duration

	// 接收, 发送处理器
	recvHandler  []cellnet.EventHandler
	sendHandler  []cellnet.EventHandler
	handlerGuard sync.RWMutex

	// 运行状态
	running      bool
	runningGuard sync.RWMutex

	// 停止过程同步
	stopping chan bool

	// 自带派发器
	*cellnet.DispatcherHandler

	// 自定义流
	streamGen func(net.Conn) cellnet.PacketStream
}

func (self *peerBase) waitStopFinished() {
	// 如果正在停止时, 等待停止完成
	if self.stopping != nil {
		<-self.stopping
		self.stopping = nil
	}
}

func (self *peerBase) isStopping() bool {
	return self.stopping != nil
}

func (self *peerBase) startStopping() {
	self.stopping = make(chan bool)
}

func (self *peerBase) endStopping() {
	select {
	case self.stopping <- true:

	default:
		self.stopping = nil
	}
}

func (self *peerBase) IsRunning() bool {

	self.runningGuard.RLock()
	defer self.runningGuard.RUnlock()

	return self.running
}

func (self *peerBase) SetRunning(v bool) {
	self.runningGuard.Lock()
	self.running = v
	self.runningGuard.Unlock()
}

func (self *peerBase) applyConnOption(conn net.Conn) {

	if cc, ok := conn.(*net.TCPConn); ok {

		if self.connReadBuffer >= 0 {
			cc.SetReadBuffer(self.connReadBuffer)
		}

		if self.connWriteBuffer >= 0 {
			cc.SetWriteBuffer(self.connWriteBuffer)
		}

		cc.SetNoDelay(self.connNoDelay)
	}

}

func (self *peerBase) SetSocketDeadline(read, write time.Duration) {
	self.connReadTimeout = read
	self.connWriteTimeout = write
}

func (self *peerBase) SocketDeadline() (read, write time.Duration) {
	return self.connReadTimeout, self.connWriteTimeout
}

func (self *peerBase) SetSocketOption(readBufferSize, writeBufferSize int, nodelay bool) {

	self.connReadBuffer = readBufferSize
	self.connWriteBuffer = writeBufferSize
	self.connNoDelay = nodelay
}

func (self *peerBase) SetPacketStreamGenerator(callback func(net.Conn) cellnet.PacketStream) {

	self.streamGen = callback
}

func (self *peerBase) genPacketStream(conn net.Conn) cellnet.PacketStream {

	self.applyConnOption(conn)

	if self.streamGen == nil {
		return NewTLVStream(conn)
	}

	return self.streamGen(conn)
}

func (self *peerBase) Queue() cellnet.EventQueue {
	return self.EventQueue
}

func (self *peerBase) nameOrAddress() string {
	if self.name != "" {
		return self.name
	}

	return self.address
}

func (self *peerBase) Tag() interface{} {
	return self.tag
}

func (self *peerBase) SetTag(tag interface{}) {
	self.tag = tag
}

func (self *peerBase) Address() string {
	return self.address
}

func (self *peerBase) SetAddress(address string) {
	self.address = address
}

func (self *peerBase) SetHandlerList(recv, send []cellnet.EventHandler) {
	self.handlerGuard.Lock()
	self.recvHandler = recv
	self.sendHandler = send
	self.handlerGuard.Unlock()
}

func (self *peerBase) HandlerList() (recv, send []cellnet.EventHandler) {
	self.handlerGuard.RLock()
	recv = self.recvHandler
	send = self.sendHandler
	self.handlerGuard.RUnlock()

	return
}

func (self *peerBase) safeRecvHandler() (ret []cellnet.EventHandler) {
	self.handlerGuard.RLock()
	ret = self.recvHandler
	self.handlerGuard.RUnlock()

	return
}

func (self *peerBase) SetName(name string) {
	self.name = name
}

func (self *peerBase) Name() string {
	return self.name
}

func (self *peerBase) SetMaxPacketSize(size int) {
	self.maxPacketSize = size
}

func (self *peerBase) MaxPacketSize() int {
	return self.maxPacketSize
}

func newPeerBase(queue cellnet.EventQueue) *peerBase {

	self := &peerBase{
		EventQueue:        queue,
		DispatcherHandler: cellnet.NewDispatcherHandler(),
		connWriteBuffer:   -1,
		connReadBuffer:    -1,
	}

	self.recvHandler = BuildRecvHandler(self.DispatcherHandler)

	self.sendHandler = BuildSendHandler()

	return self
}
