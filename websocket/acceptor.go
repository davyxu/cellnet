package websocket

import (
	"net/http"
	"net/url"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"github.com/gorilla/websocket"
)

type wsAcceptor struct {
	cellnet.EventQueue
	// 会话管理器
	cellnet.SessionManager

	// 共享配置
	*cellnet.BasePeerImplement

	// 自带派发器
	*cellnet.DispatcherHandler
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func (self *wsAcceptor) Queue() cellnet.EventQueue {
	return self.EventQueue
}

func (self *wsAcceptor) Start(address string) cellnet.Peer {

	if self.IsRunning() {
		return self
	}

	self.SetRunning(true)

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

	http.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		ses := newSession(c, self)

		// 添加到管理器
		self.Add(ses)

		// 断开后从管理器移除
		ses.OnClose = func() {
			self.Remove(ses)
		}

		ses.run()

		// 通知逻辑
		extend.PostSystemEvent(ses, cellnet.Event_Accepted, self.SafeRecvHandler(), cellnet.Result_OK)

	})

	go func() {

		err = http.ListenAndServe(url.Host, nil)

		if err != nil {
			log.Errorln(err)
		}

	}()

	return self
}

func (self *wsAcceptor) Stop() {
	if !self.IsRunning() {
		return
	}
}

func NewAcceptor(q cellnet.EventQueue) cellnet.Peer {

	self := &wsAcceptor{
		EventQueue:        q,
		DispatcherHandler: cellnet.NewDispatcherHandler(),
		SessionManager:    cellnet.NewSessionManager(),
		BasePeerImplement: cellnet.NewBasePeer(),
	}

	self.BasePeerImplement.SetHandlerList(BuildRecvHandler(self.DispatcherHandler), BuildSendHandler())

	return self
}
