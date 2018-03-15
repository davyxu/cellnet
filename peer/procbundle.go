package peer

import (
	"errors"
	"github.com/davyxu/cellnet"
)

type MessagePoster interface {
	PostEvent(ev cellnet.Event)
}

type CoreProcBundle struct {
	proc    cellnet.MessageProcessor
	hooker  cellnet.EventHooker
	handler cellnet.EventHandler
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetEventProcessor(v cellnet.MessageProcessor) {
	self.proc = v
}

func (self *CoreProcBundle) SetEventHooker(v cellnet.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetEventHandler(v cellnet.EventHandler) {
	self.handler = v
}

var notHandled = errors.New("msg not handled")

func (self *CoreProcBundle) ReadMessage(ses cellnet.Session) (msg interface{}, err error) {

	if self.proc != nil {
		return self.proc.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.proc != nil && ev != nil {
		self.proc.OnSendMessage(ev.Session(), ev.Message())
	}
}

func (self *CoreProcBundle) PostEvent(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.handler != nil && ev != nil {
		self.handler.OnEvent(ev)
	}
}
