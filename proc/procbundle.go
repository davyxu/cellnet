package proc

import (
	"github.com/davyxu/cellnet"
)

// 处理器设置接口，由各Peer实现
type ProcessorBundle interface {

	// 设置 传输器，负责收发消息
	SetTransmitter(v cellnet.MessageTransmitter)

	// 设置 接收后，发送前的事件处理流程
	SetHooker(v cellnet.EventHooker)

	// 设置 接收后最终处理回调
	SetCallback(v cellnet.EventCallback)
}

// 让EventCallback保证放在ses的队列里，而不是并发的
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

		if input == nil {
			break
		}
	}

	return input
}

func (self MultiHooker) OnOutboundEvent(input cellnet.Event) (output cellnet.Event) {

	for _, h := range self {

		input = h.OnOutboundEvent(input)

		if input == nil {
			break
		}
	}

	return input
}

func NewMultiHooker(h ...cellnet.EventHooker) cellnet.EventHooker {

	return MultiHooker(h)
}
