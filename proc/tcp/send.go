package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/util"
	"net"
)

// 发送Length-Type-Value格式的封包流程
func SendLTVPacket(ses cellnet.Session, data interface{}) cellnet.EventResult {

	// 取Socket连接
	conn, ok := ses.Raw().(net.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	var msgData []byte
	var msgID int

	switch m := data.(type) {
	case comm.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error
		var meta *cellnet.MessageMeta
		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = cellnet.EncodeMessage(data)

		if err != nil {
			log.Errorf("send message encode error: %s", err)
			return err
		}

		msgID = meta.ID
	}

	return rawSend(conn, msgData, msgID)
}

func rawSend(conn net.Conn, msgData []byte, msgid int) cellnet.EventResult {
	// 创建封包写入器
	var pktWriter util.BinaryWriter

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(msgid)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(msgData); err != nil {
		return err
	}

	// 发送长度定界的变长封包
	return util.SendVariableLengthPacket(conn, pktWriter)
}
