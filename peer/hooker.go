package cellpeer

import (
	cellevent "github.com/davyxu/cellnet/event"
)

type InboundFunc func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg)
type OutboundFunc func(input *cellevent.SendMsg) (output *cellevent.SendMsg)

type Hooker struct {
	OnInbound  InboundFunc  // 事件传入
	OnOutbound OutboundFunc // 事件传出
}

func (self *Hooker) ProcEvent(ev *cellevent.RecvMsg) {

	if self.OnInbound != nil {
		ev = self.OnInbound(ev)
	}
}
