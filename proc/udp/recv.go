package udp

import (
	"encoding/binary"
	"github.com/davyxu/cellnet/codec"
)

const MTU = 1472

func recvPacket(pktData []byte) (msg interface{}, err error) {

	// 小于包头，使用nc指令测试时，为1
	if len(pktData) < 2 {
		return nil, nil
	}

	// 用小端格式读取Size
	datasize := binary.LittleEndian.Uint16(pktData)

	// 出错，等待下次数据
	if int(datasize) != len(pktData) || datasize > MTU {
		return nil, nil
	}

	// 读取消息ID
	msgid := binary.LittleEndian.Uint16(pktData[2:])

	msgData := pktData[headerSize:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
