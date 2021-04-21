package cellpeer

import (
	cellevent "github.com/davyxu/cellnet/event"
)

type Hooker struct {
	Inbound  func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent)
	Outbound func(input *cellevent.SendMsgEvent) (output *cellevent.SendMsgEvent)
}

func (self *Hooker) ProcEvent(ev *cellevent.RecvMsgEvent) {

	if self.Inbound != nil {
		ev = self.Inbound(ev)
	}
}
