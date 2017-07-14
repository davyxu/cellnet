package websocket

import (
	"github.com/davyxu/cellnet"
	"github.com/gorilla/websocket"
)

type ReadPacketHandler struct {
}

func errToResult(err error) cellnet.Result {

	if err == nil {
		return cellnet.Result_OK
	}

	// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {

	return cellnet.Result_SocketError
}

func (self *ReadPacketHandler) Call(ev *cellnet.Event) {

	switch ev.Type {
	case cellnet.Event_Recv:

		rawSes := ev.Ses.(*wsSession)

		// 读超时
		t, data, err := rawSes.conn.ReadMessage()

		if err != nil {
			ev.SetResult(errToResult(err))
			return
		}

		switch t {
		case websocket.TextMessage:

			msgName, userdata := parsePacket(data)

			if msgName != "" {

				meta := cellnet.MessageMetaByName(msgName)

				if meta == nil || meta.Codec == nil {
					ev.SetResult(cellnet.Result_CodecError)
					return
				}

				ev.MsgID = meta.ID
				ev.Data = userdata

			} else {

				ev.Data = data
			}

		case websocket.CloseMessage:
			ev.SetResult(cellnet.Result_RequestClose)
		}

	}

}

var defaultReadPacketHandler = new(ReadPacketHandler)

func StaticReadPacketHandler() cellnet.EventHandler {
	return defaultReadPacketHandler
}
