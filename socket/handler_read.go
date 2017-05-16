package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type ReadPacketHandler struct {
	cellnet.BaseEventHandler
}

func (self *ReadPacketHandler) Call(ev *cellnet.SessionEvent) {

	switch ev.Type {
	case cellnet.SessionEvent_Recv:

		rawSes := ev.Ses.(*SocketSession)

		msgid, data, err := rawSes.stream.Read()

		if err != nil {

			castToSystemEvent(ev, cellnet.SessionEvent_Closed, &coredef.SessionClosed{Reason: err.Error()})

			ev.EndRecvLoop = true
		} else {

			ev.MsgID = msgid
			// 逻辑封包
			ev.Data = data
		}

	}

}

func NewReadPacketHandler() cellnet.EventHandler {

	return &ReadPacketHandler{}

}
