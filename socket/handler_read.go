package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type ReadPacketHandler struct {
	cellnet.BaseEventHandler

	q cellnet.EventQueue
}

func (self *ReadPacketHandler) Call(ev *cellnet.SessionEvent) {

	switch ev.Type {
	case cellnet.SessionEvent_Recv:

		rawSes := ev.Ses.(*SocketSession)

		err := rawSes.stream.Read(ev)

		if err != nil {

			castToSystemEvent(ev, cellnet.SessionEvent_Closed, &coredef.SessionClosed{Reason: err.Error()})

			ev.EndRecvLoop = true
		}

		// 逻辑封包
	}

}

func NewReadPacketHandler(q cellnet.EventQueue) cellnet.EventHandler {

	return &ReadPacketHandler{
		q: q,
	}

}
