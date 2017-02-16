package cellnet

import "reflect"

type DecodePacketHandler struct {
	BaseEventHandler
	meta *MessageMeta
}

func (self *DecodePacketHandler) Call(ev *SessionEvent) (err error) {

	// 创建消息
	ev.Msg = reflect.New(self.meta.Type).Interface()

	var codec Codec
	if ev.OverrideCodec != nil {
		codec = ev.OverrideCodec
	} else {
		codec = ev.PacketCodec()
	}

	if codec == nil {
		panic("require codec")
	}

	// 解析消息
	err = codec.Decode(ev.Data, ev.Msg)

	if err != nil {
		log.Errorf("unmarshaling error: %v, raw: %v", err, ev)
		return
	}

	return self.CallNext(ev)
}

func NewDecodePacketHandler(meta *MessageMeta) EventHandler {
	return &DecodePacketHandler{
		meta: meta,
	}
}
