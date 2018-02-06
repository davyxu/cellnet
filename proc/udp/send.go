package udp

import (
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/util"
)

func SendLTVPacket(writer udp.DataWriter, msg interface{}) error {

	// 将用户数据转换为字节数组和消息ID
	data, meta, err := codec.EncodeMessage(msg)

	if err != nil {
		log.Errorf("send message encode error: %s", err)
		return err
	}

	// 创建封包写入器
	var pktWriter util.BinaryWriter

	// 写入消息长度做验证
	if err := pktWriter.WriteValue(uint16(len(data)) + 2 + 2); err != nil {
		return err
	}

	// 写入消息ID
	if err := pktWriter.WriteValue(uint16(meta.ID)); err != nil {
		return err
	}

	// 写入序列化好的消息数据
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}

	writer.WriteData(pktWriter.Raw())

	return nil
}
