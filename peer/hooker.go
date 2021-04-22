package cellpeer

import (
	cellevent "github.com/davyxu/cellnet/event"
)

type InboundFunc func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent)
type OutboundFunc func(input *cellevent.SendMsgEvent) (output *cellevent.SendMsgEvent)

type Hooker struct {
	Inbound  InboundFunc
	Outbound OutboundFunc
}

func (self *Hooker) ProcEvent(ev *cellevent.RecvMsgEvent) {

	if self.Inbound != nil {
		ev = self.Inbound(ev)
	}
}
