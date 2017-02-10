package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
)

type ReadPacketHandler struct {
	next cellnet.Handler

	q cellnet.EventQueue
}

func (self *ReadPacketHandler) Call(evid int, data interface{}) (ret error) {

	ev := data.(*SessionEvent)

	switch evid {
	case SessionEvent_Recv:

		rawSes := ev.Ses.(*SocketSession)

		pkt, err := rawSes.stream.Read()

		if err != nil {

			ev.MsgID = Event_SessionClosed
			ev.Packet, _ = cellnet.BuildPacket(&gamedef.SessionClosed{Reason: err.Error()})
			ret = err

		} else {
			ev.MsgID = pkt.MsgID
			ev.Packet = pkt
		}

		// 逻辑封包

	case SessionEvent_Accepted:
		ev.MsgID = Event_SessionClosed
	case SessionEvent_Connected:
		ev.MsgID = Event_SessionConnected
	case SessionEvent_AcceptFailed:
		ev.MsgID = Event_SessionAcceptFailed
	case SessionEvent_ConnectFailed:
		ev.MsgID = Event_SessionConnectFailed
	}

	msgLog("recv", ev.Ses, ev.Packet)

	self.q.Post(nil, func() {
		self.next.Call(SessionEvent_Recv, ev)
	})

	return
}

func NewReadPacketHandler(next cellnet.Handler, q cellnet.EventQueue) cellnet.Handler {

	return &ReadPacketHandler{
		next: next,
		q:    q,
	}

}

type WritePacketHandler struct {
}

func (self *WritePacketHandler) Call(evid int, data interface{}) error {

	return nil
}
