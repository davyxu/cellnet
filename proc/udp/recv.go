package udp

import (
	"encoding/binary"
	"github.com/davyxu/cellnet/codec"
)

const (
	MTU       = 1472 // 最大传输单元
	packetLen = 2    // 包体大小字段
	MsgIDLen  = 2    // 消息ID字段

	HeaderSize = MsgIDLen + MsgIDLen // 整个UDP包头部分
)

func RecvPacket(pktData []byte) (msg interface{}, err error) {

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

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
