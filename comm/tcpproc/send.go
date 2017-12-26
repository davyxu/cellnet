package tcpproc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"net"
)

// 发送Length-Type-Value格式的封包流程
func SendLTVPacket(ses cellnet.Session, msg interface{}) cellnet.EventResult {

	// 取Socket连接
	conn, ok := ses.Raw().(net.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	// 将用户数据转换为字节数组和消息ID
	data, meta, err := cellnet.EncodeMessage(msg)

	if err != nil {
		log.Errorf("send message encode error: %s", err)
		return err
	}

	// 创建封包写入器
	var pktWriter util.BinaryWriter

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(meta.ID)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}

	// 发送长度定界的变长封包
	return SendVariableLengthPacket(conn, pktWriter)
}
