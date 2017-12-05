package tcppkt

import (
	"github.com/davyxu/cellnet"
	"io"
	"net"
)

// 接收Length-Type-Value格式的封包流程
func onRecvLTVPacket(ses cellnet.Session, eventFunc cellnet.EventFunc) error {

	// 取Socket连接
	conn, ok := ses.Raw().(net.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	// 接收长度定界的变长封包，返回封包读取器
	pktReader, err := RecvVariableLengthPacket(conn)

	if err != nil {
		return err
	}

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
	eventFunc(RecvMsgEvent{ses, msg})

	return nil
}

func RecvLTVPacket(inputStream io.Reader) (msg interface{}, msgid uint16, err error) {

	// 接收长度定界的变长封包，返回封包读取器
	pktReader, err := RecvVariableLengthPacket(inputStream)

	if err != nil {
		return
	}

	// 读取消息ID
	if err = pktReader.ReadValue(&msgid); err != nil {
		return
	}

	msgData := pktReader.RemainBytes()

	// 将字节数组和消息ID用户解出消息
	msg, err = cellnet.DecodeMessage(uint32(msgid), msgData)

	return
}
