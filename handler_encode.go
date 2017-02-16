package cellnet

type EncodePacketHandler struct {
	BaseEventHandler
}

func (self *EncodePacketHandler) Call(ev *SessionEvent) error {

	ev.FromMessage(ev.Msg)

	return self.CallNext(ev)

}

func NewEncodePacketHandler() EventHandler {
	return &EncodePacketHandler{}
}
