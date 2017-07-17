package websocket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"github.com/gorilla/websocket"
)

type wsSession struct {
	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	conn *websocket.Conn
}

func (self *wsSession) ID() int64 {
	return self.id
}

func (self *wsSession) SetID(id int64) {
	self.id = id
}

func (self *wsSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *wsSession) Close() {

	self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (self *wsSession) Send(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Send, self)
	ev.Msg = data

	self.RawSend(ev.SendHandler, ev)

}

func (self *wsSession) Post(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Post, self)

	ev.Msg = data

	cellnet.MsgLog(ev)

	self.p.Call(ev)
}

func (self *wsSession) RawPost(recvHandler []cellnet.EventHandler, ev *cellnet.Event) {
	if recvHandler == nil {
		recvHandler, _ = self.p.HandlerList()
	}

	ev.Ses = self

	cellnet.HandlerChainCall(recvHandler, ev)
}

func (self *wsSession) RawSend(sendHandler []cellnet.EventHandler, ev *cellnet.Event) {

	if sendHandler == nil {
		_, sendHandler = self.p.HandlerList()
	}

	ev.Ses = self

	cellnet.HandlerChainCall(sendHandler, ev)
}

func (self *wsSession) recvThread() {

	recvList, _ := self.p.HandlerList()

	for {

		ev := cellnet.NewEvent(cellnet.Event_Recv, self)

		cellnet.HandlerChainCall(recvList, ev)

		if ev.Result() != cellnet.Result_OK {

			extend.PostSystemEvent(ev.Ses, cellnet.Event_Closed, recvList, ev.Result())
			break
		}
	}
}

func (self *wsSession) run() {

	go self.recvThread()
}

func newSession(c *websocket.Conn, p cellnet.Peer) *wsSession {

	self := &wsSession{
		p:    p,
		conn: c,
	}

	return self
}
