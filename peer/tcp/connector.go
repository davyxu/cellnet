package tcp

import (
	"context"
	"github.com/davyxu/cellnet/event"
	"net"
	"sync"
	"time"
)

type Connector struct {
	*Peer

	// 连接会话
	Session *Session

	// 连接超时
	ConnTimeout time.Duration

	// 重连间隔
	ReconnInterval time.Duration

	// 连接地址
	Address string

	endSignal sync.WaitGroup

	cancelFunc context.CancelFunc

	tryConnTimes int32 // 尝试连接次数
}

// 发起连接到指定地址
func (self *Connector) Connect(address string) error {
	self.Address = address

	var ctx context.Context
	ctx, self.cancelFunc = context.WithCancel(context.Background())
	return self.conn(ctx)
}

// 异步连接到指定地址
func (self *Connector) AsyncConnect(address string) {
	self.Address = address

	var ctx context.Context
	ctx, self.cancelFunc = context.WithCancel(context.Background())
	go self.conn(ctx)
}

// 关闭连接
func (self *Connector) Close() {
	self.Session.Close()

	if self.cancelFunc != nil {
		self.cancelFunc()
	}

	self.tryConnTimes = 0
}

// 连接的端口
func (self *Connector) Port() int {
	if self.Session.conn == nil {
		return 0
	}

	return self.Session.conn.LocalAddr().(*net.TCPAddr).Port
}

func (self *Connector) conn(ctx context.Context) (err error) {

	var connectedTimes int32
	for {
		self.tryConnTimes++

		d := net.Dialer{Timeout: self.ConnTimeout}
		var conn net.Conn
		conn, err = d.DialContext(ctx, "tcp", self.Address)

		if err != nil {

			// 手动关闭时, 不要重连
			if self.Session.IsManualClosed() {
				break
			}

			self.ProcEvent(cellevent.BuildSystemEvent(self.Session, &cellevent.SessionConnectError{
				Err:            err,
				RetryTimes:     self.tryConnTimes,
				ConnectedTimes: connectedTimes,
			}))

			if self.ReconnInterval == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnInterval)

			continue
		}

		self.Session.conn = conn

		self.endSignal.Add(1)

		self.ApplySocketOption(conn)

		self.Session.Start()

		self.tryConnTimes = 0
		connectedTimes++

		self.ProcEvent(cellevent.BuildSystemEvent(self.Session, &cellevent.SessionConnected{}))

		self.endSignal.Wait()

		// 连接断开了, 没重连就退出循环
		if self.ReconnInterval == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnInterval)

		continue
	}

	return
}

func NewConnector() *Connector {
	self := &Connector{
		Peer:        newPeer(),
		ConnTimeout: time.Second * 5,
	}

	self.Session = newSession(nil, self.Peer, self)
	self.Session.endNotify = func() {
		self.endSignal.Done()
	}

	self.Init()

	return self
}
