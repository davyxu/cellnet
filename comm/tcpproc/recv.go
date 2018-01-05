package tcpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"net"
)

// 接收Length-Type-Value格式的封包流程
func RecvLTVPacket(ses cellnet.Session) (msg interface{}, err error) {

	// 取Socket连接
	conn, ok := ses.Raw().(net.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil, nil
	}

	// 接收长度定界的变长封包，返回封包读取器
	pktReader, err := util.RecvVariableLengthPacket(conn)

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
	msg, _, err = cellnet.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
