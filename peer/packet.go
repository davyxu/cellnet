package cellpeer

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	"github.com/davyxu/x/container"
)

// 将用户消息封装为发送事件
func PackEvent(payload any, ps *xcontainer.Mapper) *cellevent.SendMsg {
	var (
		msgData []byte
		msgId   int
		meta    *cellmeta.Meta
	)

	switch raw := payload.(type) {
	case *RawPacket: // 发裸包
		msgData = raw.MsgData
		msgId = raw.MsgId
	default: // 发普通编码包
		// 将用户数据转换为字节数组和消息ID
		msgData, meta, _ = cellcodec.Encode(payload, ps)

		if meta != nil {
			msgId = meta.Id
		} else {
			// 无法识别的消息, 丢给transmitter层处理
			return &cellevent.SendMsg{Msg: payload}
		}

	}
	return &cellevent.SendMsg{
		MsgId:   msgId,
		MsgData: msgData,
	}
}

// 直接发送数据时，将*RawPacket作为Send参数
type RawPacket struct {
	MsgData []byte
	MsgId   int
}

func (self *RawPacket) MessageId() int {
	return self.MsgId
}

func (self *RawPacket) Message() any {

	// 获取消息元信息
	meta := cellmeta.MetaById(self.MsgId)

	// 消息没有注册
	if meta == nil {
		return struct{}{}
	}

	// 创建消息
	msg := meta.NewType()

	// 从字节数组转换为消息
	err := meta.Codec.Decode(self.MsgData, msg)
	if err != nil {
		return struct{}{}
	}

	return msg
}
