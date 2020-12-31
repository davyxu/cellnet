package tcp

import (
	"context"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/ulog"
	xframe "github.com/davyxu/x/frame"
	"net"
	"sync"
	"time"
)

type tcpConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	xframe.PropertySet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *tcpSession

	tryConnTimes int // 尝试连接次数

	sesEndSignal sync.WaitGroup

	reconInterval time.Duration

	conTimeout time.Duration

	reconReportLimitTimes int

	cancelFunc context.CancelFunc
}

func (self *tcpConnector) Start() cellnet.Peer {

	if self.IsRunning() {
		return self
	}

	var ctx context.Context
	ctx, self.cancelFunc = context.WithCancel(context.Background())

	go self.connect(self.Address(), ctx)

	return self
}

func (self *tcpConnector) Stop() {
	if !self.IsRunning() {
		return
	}

	self.StartStopping()

	// 通知发送关闭, session设置Manual close
	self.defaultSes.Close()

	if self.cancelFunc != nil {
		self.cancelFunc()
	}

	self.tryConnTimes = 0

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *tcpConnector) ReconnectInterval() time.Duration {

	return self.reconInterval
}

func (self *tcpConnector) SetReconnectInterval(v time.Duration) {
	self.reconInterval = v
}

func (self *tcpConnector) SetConnectTimeout(v time.Duration) {
	self.conTimeout = v
}

func (self *tcpConnector) SetReconnectReportLimitTimes(v int) {
	self.reconReportLimitTimes = v
}

func (self *tcpConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *tcpConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *tcpConnector) Port() int {

	conn := self.defaultSes.Conn()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.TCPAddr).Port
}

// 连接器，传入连接地址和发送封包次数
func (self *tcpConnector) connect(address string, ctx context.Context) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		d := net.Dialer{Timeout: self.conTimeout}
		// 尝试用Socket连接地址
		conn, err := d.DialContext(ctx, "tcp", address)

		self.defaultSes.setConn(conn)

		// 发生错误时退出
		if err != nil {

			// 直接关闭时，退出连接循环
			if self.defaultSes.IsManualClosed() {
				break
			}

			if self.tryConnTimes <= self.reconReportLimitTimes {
				ulog.Errorf("#tcp.connect failed(%s), times: %d %v %p", self.Name(), self.tryConnTimes, err.Error(), self)
			}

			// 没重连就退出
			if self.reconInterval == 0 {

				self.ProcEvent(cellnet.BuildSystemEvent(self.defaultSes, &cellnet.SessionConnectError{}))
				break
			}

			// 有重连就等待
			time.Sleep(self.reconInterval)

			// 继续连接
			continue
		}

		self.sesEndSignal.Add(1)

		self.ApplySocketOption(conn)

		self.defaultSes.Start()

		self.tryConnTimes = 0

		self.ProcEvent(cellnet.BuildSystemEvent(self.defaultSes, &cellnet.SessionConnected{}))

		self.sesEndSignal.Wait()

		self.defaultSes.setConn(nil)

		// 没重连就退出/主动退出
		if self.reconInterval == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.reconInterval)

		// 继续连接
		continue

	}

	self.EndStopping()
	self.SetRunning(false)
}

func (self *tcpConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *tcpConnector) TypeName() string {
	return "tcp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &tcpConnector{
			SessionManager:        new(peer.CoreSessionManager),
			reconReportLimitTimes: 3,
			conTimeout:            time.Second * 3,
		}

		self.defaultSes = newSession(nil, self, func() {
			self.sesEndSignal.Done()
		})

		self.CoreTCPSocketOption.Init()

		return self
	})
}
