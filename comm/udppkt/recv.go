package udppkt

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

var ErrPacketCrack = errors.New("udp packet crack by len")

const MTU = 1472

func onRecvLTVPacket(ses cellnet.Session, data []byte, eventFunc cellnet.EventFunc) cellnet.EventResult {

	var pktReader util.BinaryReader
	pktReader.Init(data)

	// 读取消息ID
	var datasize uint16
	if err := pktReader.ReadValue(&datasize); err != nil {
		return err
	}

	// 出错，等待下次数据
	if int(datasize) != len(data) || datasize > MTU {
		return nil
	}

	// 读取消息ID
	var msgid uint16
	if err := pktReader.ReadValue(&msgid); err != nil {
		return err
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, _, err := cellnet.DecodeMessage(uint32(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return err
	}

	// 调用用户回调
	return eventFunc(&cellnet.RecvMsgEvent{ses, msg})
}
