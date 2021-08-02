package cellpeer

import (
	cellcodec "github.com/davyxu/cellnet/codec"
	cellevent "github.com/davyxu/cellnet/event"
	cellmeta "github.com/davyxu/cellnet/meta"
	xframe "github.com/davyxu/x/frame"
)

// 将用户消息封装为发送事件
func PackEvent(payload interface{}, ps *xframe.Mapper) *cellevent.SendMsg {
	var (
		msgData []byte
		msgID   int
		meta    *cellmeta.Meta
	)

	switch raw := payload.(type) {
	case *RawPacket: // 发裸包
		msgData = raw.MsgData
		msgID = raw.MsgID
	default: // 发普通编码包
		// 将用户数据转换为字节数组和消息ID
		msgData, meta, _ = cellcodec.Encode(payload, ps)

		if meta != nil {
			msgID = meta.ID
		} else {
			// 无法识别的消息, 丢给transmitter层处理
			return &cellevent.SendMsg{Msg: payload}
		}

	}
	return &cellevent.SendMsg{
		MsgID:   msgID,
		MsgData: msgData,
	}
}

// 直接发送数据时，将*RawPacket作为Send参数
type RawPacket struct {
	MsgData []byte
	MsgID   int
}

func (self *RawPacket) MessageID() int {
	return self.MsgID
}

func (self *RawPacket) Message() interface{} {

	// 获取消息元信息
	meta := cellmeta.MetaByID(self.MsgID)

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
