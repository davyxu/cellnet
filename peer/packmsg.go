package peer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

func PackEvent(payload interface{}, ctx cellnet.ContextSet) *cellnet.SendMsgEvent {
	var (
		msgData []byte
		msgID   int
		meta    *cellnet.MessageMeta
	)

	switch raw := payload.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = raw.MsgData
		msgID = raw.MsgID
	default: // 发普通编码包
		// 将用户数据转换为字节数组和消息ID
		msgData, meta, _ = codec.EncodeMessage(payload, ctx)

		if meta != nil {
			msgID = meta.ID
		} else {
			// 无法识别的消息, 丢给transmitter层处理
			return &cellnet.SendMsgEvent{Msg: payload}
		}

	}
	return &cellnet.SendMsgEvent{
		MsgID:   msgID,
		MsgData: msgData,
	}
}
