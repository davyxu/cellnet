package tcptransport

import (
	"fmt"
	cellcodec "github.com/davyxu/cellnet/codec"
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/tcp"
	xbytes "github.com/davyxu/x/bytes"
	xio "github.com/davyxu/x/io"
)

type NoCryptMessage struct {
	Msg any
}

func (self *NoCryptMessage) Message() any {
	return self.Msg
}

func SendMessage(ses *tcp.Session, ev *cellevent.SendMsg) error {

	if TestEnableSendPanic {
		panic("emulate send crash")
	}

	ps := &ses.Peer.Mapper

	var (
		msgData []byte
	)

	if ev.MessageData() != nil {
		msgData = ev.MessageData()
	} else if raw, ok := ev.Message().(*NoCryptMessage); ok {
		data, meta, err := cellcodec.Encode(raw.Msg, ps)
		if err != nil {
			return fmt.Errorf("encode msg failed, %+v", raw.Msg)
		} else {
			msgData = data
			ev.MsgId = meta.Id
		}
	} else {
		panic(fmt.Sprintf("invalid message %+v", ev.Message()))
	}

	bodySize := msgIDLen + len(msgData)
	composeBuffer := make([]byte, packetHeaderSize+bodySize)
	writer := xbytes.NewWriter(composeBuffer)

	writer.WriteUint32(uint32(bodySize))
	writer.WriteUint32(uint32(ev.MsgId))
	writer.Write(msgData)

	// 将数据写入Socket

	return xio.WriteFull(ses.Raw(), composeBuffer)
}
