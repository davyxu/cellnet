package udptransmit

import (
	"encoding/binary"
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/udp"
)

func RecvMessage(ses *udp.Session, pktData []byte) (ev *cellevent.RecvMsg, err error) {
	// 小于包头，使用nc指令测试时，为1
	if len(pktData) < packetLen {
		return nil, nil
	}

	// 用小端格式读取Size
	datasize := binary.LittleEndian.Uint16(pktData)

	// 出错，等待下次数据
	if int(datasize) != len(pktData) || datasize > MTU {
		return nil, nil
	}

	// 读取消息ID
	msgid := binary.LittleEndian.Uint16(pktData[packetLen:])

	msgData := pktData[HeaderSize:]

	ev = &cellevent.RecvMsg{
		Ses:     ses,
		MsgID:   int(msgid),
		MsgData: msgData,
	}

	return
}
