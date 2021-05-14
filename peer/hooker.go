package cellpeer

import (
	cellevent "github.com/davyxu/cellnet/event"
)

type InboundFunc func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg)
type OutboundFunc func(input *cellevent.SendMsg) (output *cellevent.SendMsg)

type Hooker struct {
	Inbound  InboundFunc  // 事件传入
	Outbound OutboundFunc // 事件传出
}

func (self *Hooker) ProcEvent(ev *cellevent.RecvMsg) {

	if self.Inbound != nil {
		ev = self.Inbound(ev)
	}
}
