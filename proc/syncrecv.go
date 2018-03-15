package proc

import (
	"github.com/davyxu/cellnet"
)

// 直接回调用户回调

type SyncReceiver struct {
	evChan chan cellnet.Event
}

func (self *SyncReceiver) Recv(callback cellnet.EventCallback) *SyncReceiver {
	callback(<-self.evChan)
	return self
}

func NewSyncReceiver(p cellnet.Peer) *SyncReceiver {

	self := &SyncReceiver{
		evChan: make(chan cellnet.Event),
	}
	p.(interface {
		SetEventCallback(v cellnet.EventCallback)
	}).SetEventCallback(func(ev cellnet.Event) {

		self.evChan <- ev
	})

	return self
}
