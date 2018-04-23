package proc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

// 同步接收消息器, 可选件，可作为流程测试辅助工具
type SyncReceiver struct {
	evChan chan cellnet.Event

	callback func(ev cellnet.Event)
}

// 将处理回调返回给BindProcessorHandler用于注册
func (self *SyncReceiver) EventCallback() cellnet.EventCallback {

	return self.callback
}

// 持续阻塞，直到某个消息到达后，使用回调返回消息
func (self *SyncReceiver) Recv(callback cellnet.EventCallback) *SyncReceiver {
	callback(<-self.evChan)
	return self
}

// 持续阻塞，直到某个消息到达后，返回消息
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
