package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
	"time"
)

type ReadPacketHandler struct {
}

func (self *ReadPacketHandler) Call(ev *cellnet.SessionEvent) {

	switch ev.Type {
	case cellnet.SessionEvent_Recv:

		rawSes := ev.Ses.(*SocketSession)

		// 读超时
		read, _ := rawSes.FromPeer().SocketDeadline()

		if read != 0 {
			rawSes.stream.Raw().SetReadDeadline(time.Now().Add(read))
		}

		msgid, data, err := rawSes.stream.Read()

		if err != nil {

			castToSystemEvent(ev, cellnet.SessionEvent_Closed, &coredef.SessionClosed{Reason: int32(errToReason(err))})

			ev.EndRecvLoop = true
		} else {

			ev.MsgID = msgid
			// 逻辑封包
			ev.Data = data
		}

	}

}

var defaultReadPacketHandler = new(ReadPacketHandler)

func StaticReadPacketHandler() cellnet.EventHandler {
	return defaultReadPacketHandler
}
