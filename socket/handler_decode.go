package socket

import (
	"reflect"

	"github.com/davyxu/cellnet"
)

type DecodePacketHandler struct {
	cellnet.BaseEventHandler
	meta *cellnet.MessageMeta
}

func (self *DecodePacketHandler) Call(ev *cellnet.SessionEvent) (err error) {

	if ev.Data == nil {
		ev.Msg = reflect.New(self.meta.Type.Elem()).Interface()
	} else {

		ev.Meta = self.meta

		ev.Msg, err = cellnet.ParsePacket(ev.Data, self.meta.Type)

		if err != nil {
			log.Errorf("unmarshaling error: %v, raw: %v", err, ev)
			return
		}
	}

	return self.CallNext(ev)
}

func NewDecodePacketHandler(meta *cellnet.MessageMeta) cellnet.EventHandler {
	return &DecodePacketHandler{
		meta: meta,
	}
}
