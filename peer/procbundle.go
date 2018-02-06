package peer

import (
	"errors"
	"github.com/davyxu/cellnet"
)

type MessagePoster interface {
	PostEvent(ev cellnet.Event)
}

type CoreProcessorBundle struct {
	proc    cellnet.MessageProcessor
	hooker  cellnet.EventHooker
	handler cellnet.EventHandler
}

func (self *CoreProcessorBundle) GetBundle() *CoreProcessorBundle {
	return self
}

func (self *CoreProcessorBundle) SetEventProcessor(v cellnet.MessageProcessor) {
	self.proc = v
}

func (self *CoreProcessorBundle) SetEventHooker(v cellnet.EventHooker) {
	self.hooker = v
}

func (self *CoreProcessorBundle) SetEventHandler(v cellnet.EventHandler) {
	self.handler = v
}

var notHandled = errors.New("not handled")

func (self *CoreProcessorBundle) ReadMessage(ses cellnet.Session) (msg interface{}, err error) {

	if self.proc != nil {
		return self.proc.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcessorBundle) SendMessage(ev cellnet.Event) {

	if self.hooker != nil {
		self.hooker.OnOutboundEvent(ev)
	}

	if self.proc != nil {
		self.proc.OnSendMessage(ev.BaseSession(), ev.Message())
	}
}

func (self *CoreProcessorBundle) PostEvent(ev cellnet.Event) {

	if self.hooker != nil {
		self.hooker.OnInboundEvent(ev)
	}

	if self.handler != nil {
		self.handler.OnEvent(ev)
	}
}
