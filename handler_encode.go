package cellnet

type EncodePacketHandler struct {
	BaseEventHandler
}

func (self *EncodePacketHandler) Call(ev *SessionEvent) {

	ev.FromMessage(ev.Msg)

}

func NewEncodePacketHandler() EventHandler {
	return &EncodePacketHandler{}
}
