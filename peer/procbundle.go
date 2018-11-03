package peer

import (
	"errors"
	"github.com/davyxu/cellnet"
)

// 手动投递消息， 兼容v2的设计
type MessagePoster interface {

	// 投递一个消息到Hooker之前
	ProcEvent(ev cellnet.Event)
}

type CoreProcBundle struct {
	transmit cellnet.MessageTransmitter
	hooker   cellnet.EventHooker
	callback cellnet.EventCallback
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetTransmitter(v cellnet.MessageTransmitter) {
	self.transmit = v
}

func (self *CoreProcBundle) SetHooker(v cellnet.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetCallback(v cellnet.EventCallback) {
	self.callback = v
}

var notHandled = errors.New("Processor: Transimitter nil")

func (self *CoreProcBundle) ReadMessage(ses cellnet.Session) (msg interface{}, err error) {

	if self.transmit != nil {
		return self.transmit.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.transmit != nil && ev != nil {
		self.transmit.OnSendMessage(ev.Session(), ev.Message())
	}
}

func (self *CoreProcBundle) ProcEvent(ev cellnet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.callback != nil && ev != nil {
		self.callback(ev)
	}
}
