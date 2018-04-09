package proc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

// 直接回调用户回调

type SyncReceiver struct {
	evChan chan cellnet.Event

	callback func(ev cellnet.Event)
}

func (self *SyncReceiver) EventCallback() cellnet.EventCallback {

	return self.callback
}

func (self *SyncReceiver) Recv(callback cellnet.EventCallback) *SyncReceiver {
	callback(<-self.evChan)
	return self
}

func (self *SyncReceiver) WaitMessage(msgName string) (msg interface{}) {

	var wg sync.WaitGroup

	meta := cellnet.MessageMetaByFullName(msgName)
	if meta == nil {
		panic("unknown message name:" + msgName)
	}

	wg.Add(1)

	self.Recv(func(ev cellnet.Event) {

		inMeta := cellnet.MessageMetaByType(reflect.TypeOf(ev.Message()))
		if inMeta == meta {
			msg = ev.Message()
			wg.Done()
		}

	})

	wg.Wait()
	return
}

func NewSyncReceiver(p cellnet.Peer) *SyncReceiver {

	self := &SyncReceiver{
		evChan: make(chan cellnet.Event),
	}

	self.callback = func(ev cellnet.Event) {

		self.evChan <- ev
	}

	return self
}
