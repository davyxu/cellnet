package udppkt

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

func onRecvTVPacket(ses cellnet.Session, data []byte, eventFunc cellnet.EventFunc) cellnet.EventResult {

	var pktReader util.BinaryReader
	pktReader.Init(data)

	// 读取消息ID
	var msgid uint16
	if err := pktReader.ReadValue(&msgid); err != nil {
		return err
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, err := cellnet.DecodeMessage(uint32(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return err
	}

	// 调用用户回调
	return eventFunc(cellnet.RecvMsgEvent{ses, msg})
}
