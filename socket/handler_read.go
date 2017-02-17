package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/gamedef"
)

type ReadPacketHandler struct {
	cellnet.BaseEventHandler

	q cellnet.EventQueue
}

func (self *ReadPacketHandler) Call(ev *cellnet.SessionEvent) (ret error) {

	switch ev.Type {
	case cellnet.SessionEvent_Recv:

		rawSes := ev.Ses.(*SocketSession)

		err := rawSes.stream.Read(ev)

		if err != nil {

			castToSystemEvent(ev, cellnet.SessionEvent_Closed, &gamedef.SessionClosed{Reason: err.Error()})

			ret = err
		}

		// 逻辑封包
	}

	self.q.Post(func() {
		self.CallNext(ev)
	})

	return
}

func NewReadPacketHandler(q cellnet.EventQueue) cellnet.EventHandler {

	return &ReadPacketHandler{
		q: q,
	}

}
