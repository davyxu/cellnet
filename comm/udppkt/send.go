package udppkt

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

func onSendLTVPacket(ses cellnet.Session, msg interface{}) cellnet.EventResult {

	// 将用户数据转换为字节数组和消息ID
	data, msgid, err := cellnet.EncodeMessage(msg)

	if err != nil {
		return err
	}

	// 创建封包写入器
	var pktWriter util.BinaryWriter

	// 写入消息长度做验证
	if err := pktWriter.WriteValue(uint16(len(data)) + 2 + 2); err != nil {
		return err
	}

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(msgid)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}

	writer := ses.(interface {
		WriteData(data []byte) error
	})

	return writer.WriteData(pktWriter.Raw())
}
