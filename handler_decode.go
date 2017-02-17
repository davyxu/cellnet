package cellnet

import "reflect"

type DecodePacketHandler struct {
	BaseEventHandler
	meta *MessageMeta
}

func (self *DecodePacketHandler) Call(ev *SessionEvent) (err error) {

	// 系统消息不做处理
	if !ev.IsSystemEvent() {

		// 创建消息
		ev.Msg = reflect.New(self.meta.Type).Interface()

		// 解析消息
		err = ev.PacketCodec().Decode(ev.Data, ev.Msg)

		if err != nil {
			log.Errorf("unmarshaling error: %v, raw: %v", err, ev)
			return
		}
	}

	return self.CallNext(ev)
}

func NewDecodePacketHandler(meta *MessageMeta) EventHandler {
	return &DecodePacketHandler{
		meta: meta,
	}
}
