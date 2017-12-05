package packet

import (
	"github.com/davyxu/cellnet"
	"net"
)

// 发送Length-Type-Value格式的封包流程
func onSendLTVPacket(ses cellnet.Session, f cellnet.EventFunc, msg interface{}) error {

	// 取Socket连接
	conn, ok := ses.Raw().(net.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	// 调用用户回调
	invokeMsgFunc(ses, f, SendMsgEvent{ses, msg})

	// 将用户数据转换为字节数组和消息ID
	data, msgid, err := cellnet.EncodeMessage(msg)

	if err != nil {
		return err
	}

	// 创建封包写入器
	var pktWriter PacketWriter

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(msgid)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}

	// 发送长度定界的变长封包
	return SendVariableLengthPacket(conn, pktWriter)
}
