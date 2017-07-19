package websocket

import (
	"net/url"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"github.com/gorilla/websocket"
)

// 连接器, 可由Peer转换
type Connector interface {

	// 连接后的Session
	DefaultSession() cellnet.Session
}

type wsConnector struct {
	*wsPeer

	closeSignal chan bool

	defaultSes cellnet.Session
}

func (self *wsConnector) Start(address string) cellnet.Peer {

	if self.IsRunning() {
		return self
	}

	url, err := url.Parse(address)

	if err != nil {
		log.Errorln(err, address)
		return self
	}

	if url.Path == "" {
		log.Errorln("websocket: expect path in url to listen", address)
		return self
	}

	self.SetAddress(address)

	go self.connect()

	return self
}

func errToResult(err error) cellnet.Result {

	if err == nil {
		return cellnet.Result_OK
	}

	return cellnet.Result_SocketError
}

func (self *wsConnector) connect() {

	self.SetRunning(true)

	defer self.SetRunning(false)

	c, _, err := websocket.DefaultDialer.Dial(self.Address(), nil)
	if err != nil {
		extend.PostSystemEvent(nil, cellnet.Event_ConnectFailed, self.ChainListRecv(), errToResult(err))
		return
	}

	ses := newSession(c, self)

	// 添加到管理器
	self.Add(ses)

	// 断开后从管理器移除
	ses.OnClose = func() {
		self.Remove(ses)
		self.closeSignal <- true
	}

	ses.run()

	// 通知逻辑
	extend.PostSystemEvent(ses, cellnet.Event_Connected, self.ChainListRecv(), cellnet.Result_OK)

	if <-self.closeSignal {

		self.defaultSes = nil
	}
}

func (self *wsConnector) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}
}

func (self *wsConnector) DefaultSession() cellnet.Session {
	return self.defaultSes
}

func NewConnector(q cellnet.EventQueue) cellnet.Peer {

	self := &wsConnector{
		wsPeer:      newPeer(q, cellnet.NewSessionManager()),
		closeSignal: make(chan bool),
	}

	return self
}
