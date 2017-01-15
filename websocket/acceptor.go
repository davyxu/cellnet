package websocket

import (
	"net"
	"net/http"
	"strings"

	"github.com/davyxu/cellnet"
	"golang.org/x/net/websocket"
)

type socketAcceptor struct {
	*peerBase
	*sessionMgr

	listener net.Listener

	running bool
}

func (self *socketAcceptor) websocketHandler(ws *websocket.Conn) {
	conn := ws

	// 处理连接进入独立线程, 防止accept无法响应
	ses := newSession(NewPacketStream(conn), self, self)

	// 添加到管理器
	self.sessionMgr.Add(ses)

	// 断开后从管理器移除
	ses.OnClose = func() {
		self.sessionMgr.Remove(ses)
	}

	log.Infof("#accepted(%s) sid: %d", self.name, ses.ID())

	// 通知逻辑
	self.Post(self, NewSessionEvent(Event_SessionAccepted, ses, nil))
	self.Wait()
}

func (self *socketAcceptor) Start(address string) cellnet.Peer {
	//拆分下参数

	args := strings.Split(address, ",")
	address, orign := "", ""
	if len(args) >= 1 {
		address = args[0]
	}
	if len(args) >= 2 {
		orign = args[1]
	}

	http.Handle(orign, websocket.Handler(self.websocketHandler))
	//开启http服务器监听websocket
	log.Infoln("opening")
	go func() {
		self.running = true
		log.Infof("#listen(%s) %s ", self.name, address)

		err := http.ListenAndServe(address, nil)
		log.Infoln("opened")
		//	ln, err := net.Listen("tcp", address)

		//	self.listener = ln

		if err != nil {
			log.Errorf("#listen failed(%s) %v", self.name, err.Error())
			self.running = false
		}

	}()

	// 原接收线程移动到http.handler里

	return self
}

func (self *socketAcceptor) Stop() {

	if !self.running {
		return
	}

	self.running = false

	self.listener.Close()
}

func NewAcceptor(evq cellnet.EventQueue) cellnet.Peer {

	self := &socketAcceptor{
		sessionMgr: newSessionManager(),
		peerBase:   newPeerBase(evq),
	}

	return self
}
