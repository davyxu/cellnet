package cellnet

import (
	"reflect"
)

type DecodePacketHandler struct {
	BaseEventHandler
}

func (self *DecodePacketHandler) Call(ev *SessionEvent) {

	// 系统消息不做处理
	if !ev.IsSystemEvent() {

		meta := MessageMetaByID(ev.MsgID)

		if meta.Codec == nil {
			return
		}

		// 创建消息
		ev.Msg = reflect.New(meta.Type).Interface()

		// 解析消息
		meta.Codec.Decode(ev.Data, ev.Msg)
	}

}

func NewDecodePacketHandler() EventHandler {

	return &DecodePacketHandler{}
}
