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

	tag interface{}
}

func (self *wsSession) Tag() interface{} {
	return self.tag
}
func (self *wsSession) SetTag(tag interface{}) {
	self.tag = tag
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

	if ev.ChainSend != nil {
		ev.ChainSend = self.p.ChainSend()
	}

	self.RawSend(ev)

}

func (self *wsSession) RawSend(ev *cellnet.Event) {

	ev.Ses = self

	if ev.ChainSend != nil {
		ev.ChainSend.Call(ev)
	}

	// 发送日志
	cellnet.MsgLog(ev)

	go func() {

		meta := cellnet.MessageMetaByID(ev.MsgID)

		if meta == nil {
			ev.SetResult(cellnet.Result_CodecError)
			return
		}

		raw := composePacket(meta.Name, ev.Data)

		self.conn.WriteMessage(websocket.TextMessage, raw)

	}()
}

func (self *wsSession) ReadPacket() (msgid uint32, data []byte, result cellnet.Result) {

	// 读超时
	t, raw, err := self.conn.ReadMessage()

	if err != nil {
		return 0, nil, errToResult(err)
	}

	switch t {
	case websocket.TextMessage:

		msgName, userdata := parsePacket(raw)

		data = userdata

		if msgName != "" {

			meta := cellnet.MessageMetaByName(msgName)

			if meta == nil || meta.Codec == nil {
				return 0, nil, cellnet.Result_CodecError
			}

			msgid = meta.ID

		}

	case websocket.CloseMessage:
		return 0, nil, cellnet.Result_RequestClose
	}

	return msgid, data, cellnet.Result_OK
}

func (self *wsSession) recvThread() {

	for {

		msgid, data, result := self.ReadPacket()

		chainList := self.p.ChainListRecv()

		if result != cellnet.Result_OK {

			extend.PostSystemEvent(self, cellnet.Event_Closed, chainList, result)
			break

		}

		ev := cellnet.NewEvent(cellnet.Event_Recv, self)
		ev.MsgID = msgid
		ev.Data = data

		// 接收日志
		cellnet.MsgLog(ev)

		chainList.Call(ev)

		if ev.Result() != cellnet.Result_OK {
			extend.PostSystemEvent(ev.Ses, cellnet.Event_Closed, chainList, ev.Result())
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
