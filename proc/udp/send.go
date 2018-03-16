package udp

import (
	"encoding/binary"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/peer/udp"
)

func SendLTVPacket(writer udp.DataWriter, msg interface{}) error {

	// 将用户数据转换为字节数组和消息ID
	msgData, meta, err := codec.EncodeMessage(msg)

	if err != nil {
		log.Errorf("send message encode error: %s", err)
		return err
	}

	pkt := make([]byte, 2+2+len(msgData))

	// 写入消息长度做验证
	binary.LittleEndian.PutUint16(pkt, uint16(2+2+len(msgData)))

	// Type
	binary.LittleEndian.PutUint16(pkt[2:], uint16(meta.ID))

	// Value
	copy(pkt[2+2:], msgData)

	writer.WriteData(pkt)

	return nil
}
