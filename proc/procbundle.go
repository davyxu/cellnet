package proc

import (
	"github.com/davyxu/cellnet"
)

type ProcessorBundle interface {
	SetTransmitter(v cellnet.MessageTransmitter)

	SetHooker(v cellnet.EventHooker)

	SetCallback(v cellnet.EventCallback)
}

// 放队列中回调
func NewQueuedEventCallback(callback cellnet.EventCallback) cellnet.EventCallback {

	return func(ev cellnet.Event) {
		if callback != nil {
			cellnet.SessionQueuedCall(ev.Session(), func() {

				callback(ev)
			})
		}
	}
}

// 当需要多个Hooker时，使用NewMultiHooker将多个hooker合并成1个hooker处理
type MultiHooker []cellnet.EventHooker

func (self MultiHooker) OnInboundEvent(input cellnet.Event) (output cellnet.Event) {

	for _, h := range self {
		input = h.OnInboundEvent(input)
	}

	return input
}

func (self MultiHooker) OnOutboundEvent(input cellnet.Event) (output cellnet.Event) {

	for _, h := range self {
		input = h.OnOutboundEvent(input)
	}

	return input
}

func NewMultiHooker(h ...cellnet.EventHooker) cellnet.EventHooker {

	return MultiHooker(h)
}
