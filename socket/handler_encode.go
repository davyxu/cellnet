package socket

import "github.com/davyxu/cellnet"

type EncodePacketHandler struct {
	cellnet.BaseEventHandler
}

func (self *EncodePacketHandler) Call(ev *cellnet.SessionEvent) error {

	ev.FromMessage(ev.Msg)

	return self.CallNext(ev)

}

func NewEncodePacketHandler() cellnet.EventHandler {
	return &EncodePacketHandler{}
}
