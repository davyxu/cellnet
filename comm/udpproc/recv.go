package udpproc

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

var ErrPacketCrack = errors.New("udp packet crack by len")

const MTU = 1472

func RecvLTVPacket(data []byte) (msg interface{}, err error) {

	var pktReader util.BinaryReader
	pktReader.Init(data)

	// 读取数据大小，看是否完整
	var datasize uint16
	if err := pktReader.ReadValue(&datasize); err != nil {
		return nil, err
	}

	// 出错，等待下次数据
	if int(datasize) != len(data) || datasize > MTU {
		return nil, nil
	}

	// 读取消息ID
	var msgid uint16
	if err := pktReader.ReadValue(&msgid); err != nil {
		return nil, err
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, _, err = cellnet.DecodeMessage(uint32(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
