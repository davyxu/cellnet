package proc

import (
	"github.com/davyxu/cellnet"
)

type DuplexEventInitor interface {
	SetEventFunc(processor string, inboundEvent, outboundEvent cellnet.EventProc)
	SetRaw(inboundEvent, outboundEvent cellnet.EventProc)
}

type DuplexEventInvoker interface {
	CallInboundProc(ev interface{}) interface{}

	CallOutboundProc(ev interface{}) interface{}
}

type CoreDuplexEventProc struct {
	InboundProc  cellnet.EventProc
	OutboundProc cellnet.EventProc
}

func (self *CoreDuplexEventProc) SetEventFunc(processor string, inboundEvent, outboundEvent cellnet.EventProc) {
	self.InboundProc, self.OutboundProc = MakeEventProcessor(processor, inboundEvent, outboundEvent)
}

func (self *CoreDuplexEventProc) SetRaw(inboundEvent, outboundEvent cellnet.EventProc) {
	self.InboundProc, self.OutboundProc = inboundEvent, outboundEvent
}

// socket包内部派发事件
func (self *CoreDuplexEventProc) CallInboundProc(ev interface{}) interface{} {

	if self.InboundProc == nil {
		return nil
	}

	//log.Debugf("<Inbound> %T|%+v", ev, ev)

	return self.InboundProc(ev)
}

// socket包内部派发事件
func (self *CoreDuplexEventProc) CallOutboundProc(ev interface{}) interface{} {

	if self.OutboundProc == nil {
		return nil
	}

	//log.Debugf("<Outbound> %T|%+v", ev, ev)

	return self.OutboundProc(ev)
}
