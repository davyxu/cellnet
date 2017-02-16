package cellnet

import "reflect"

type DecodePacketHandler struct {
	BaseEventHandler
	meta *MessageMeta
}

func (self *DecodePacketHandler) Call(ev *SessionEvent) (err error) {

	ev.Msg = reflect.New(self.meta.Type.Elem()).Interface()

	ev.Meta = self.meta

	codec := ev.Ses.FromPeer().PacketCodec()

	if codec == nil {
		panic("require codec")
	}

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
