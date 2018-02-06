package tcp

import (
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"io"
)

// 接收Length-Type-Value格式的封包流程
func RecvLTVPacket(reader io.Reader) (msg interface{}, err error) {

	// 取Socket连接

	// 接收长度定界的变长封包，返回封包读取器
	pktReader, err := util.RecvVariableLengthPacket(reader)

	if err != nil {
		return nil, err
	}

	// 读取消息ID
	var msgid uint16
	if err := pktReader.ReadValue(&msgid); err != nil {
		return nil, err
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
