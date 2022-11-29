package udptransport

import (
	"encoding/binary"
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/udp"
)

func SendMessage(ses *udp.Session, ev *cellevent.SendMsg) error {

	pktData := make([]byte, HeaderSize+len(ev.MsgData))

	// 写入消息长度做验证
	binary.LittleEndian.PutUint16(pktData, uint16(HeaderSize+len(ev.MsgData)))

	// Type
	binary.LittleEndian.PutUint16(pktData[2:], uint16(ev.MsgId))

	// Value
	copy(pktData[HeaderSize:], ev.MsgData)

	ses.Write(pktData)

	return nil
}
