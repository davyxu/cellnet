package udp

import (
	"encoding/binary"
	cellevent "github.com/davyxu/cellnet/event"
)

const (
	MTU       = 1472 // 最大传输单元
	packetLen = 2    // 包体大小字段
	MsgIDLen  = 2    // 消息ID字段

	HeaderSize = MsgIDLen + MsgIDLen // 整个UDP包头部分
)

func RecvMessage(ses *Session, pktData []byte) (ev *cellevent.RecvMsgEvent, err error) {
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

	ev = &cellevent.RecvMsgEvent{
		Ses:     ses,
		MsgID:   int(msgid),
		MsgData: msgData,
	}

	return
}

func SendMessage(ses *Session, ev *cellevent.SendMsgEvent) error {

	pktData := make([]byte, HeaderSize+len(ev.MsgData))

	// 写入消息长度做验证
	binary.LittleEndian.PutUint16(pktData, uint16(HeaderSize+len(ev.MsgData)))

	// Type
	binary.LittleEndian.PutUint16(pktData[2:], uint16(ev.MsgID))

	// Value
	copy(pktData[HeaderSize:], ev.MsgData)

	ses.Write(pktData)

	return nil
}
