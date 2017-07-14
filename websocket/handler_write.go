package websocket

import (
	"github.com/davyxu/cellnet"
	"github.com/gorilla/websocket"
)

type WritePacketHandler struct {
}

func (self *WritePacketHandler) Call(ev *cellnet.Event) {

	go func() {
		rawSes := ev.Ses.(*wsSession)

		meta := cellnet.MessageMetaByID(ev.MsgID)

		if meta == nil {
			ev.SetResult(cellnet.Result_CodecError)
			return
		}

		raw := composePacket(meta.Name, ev.Data)

		rawSes.conn.WriteMessage(websocket.TextMessage, raw)

	}()
}

var defaultWritePacketHandler = new(WritePacketHandler)

func StaticWritePacketHandler() cellnet.EventHandler {
	return defaultWritePacketHandler
}
