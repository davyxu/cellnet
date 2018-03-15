package util

import (
	"github.com/davyxu/cellnet"
	"sync"
)

// 直接回调用户回调
type msgWaiter struct {
	msgChan chan cellnet.Event
}

func (self *msgWaiter) OnEvent(ev cellnet.Event) {

	self.msgChan <- ev
}

func (self *msgWaiter) WaitEvent() cellnet.Event {
	return <-self.msgChan
}

func SyncRecvEvent(p cellnet.Peer, onRecv func(ev cellnet.Event)) {

	var wg sync.WaitGroup

	wg.Add(1)

	setter := p.(interface {
		SetEventHandler(v cellnet.EventHandler)
	})

	w := &msgWaiter{
		msgChan: make(chan cellnet.Event),
	}

	setter.SetEventHandler(w)

	ev := w.WaitEvent()
	onRecv(ev)
}
