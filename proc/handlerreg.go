package proc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

type MessageDispatcher struct {
	handlerByType sync.Map
}

func (self *MessageDispatcher) OnEvent(ev cellnet.Event) {

	msgType := reflect.TypeOf(ev.Message())

	if handlers, ok := self.handlerByType.Load(msgType.Elem()); ok {

		for _, callback := range handlers.([]cellnet.EventCallback) {

			callback(ev)
		}

	}
}

func (self *MessageDispatcher) RegisterMessage(peer cellnet.Peer, msgName string, userCallback cellnet.EventCallback) {
	meta := cellnet.MessageMetaByFullName(msgName)
	if meta == nil {
		panic("message not found:" + msgName)
	}

	rawhandlers, _ := self.handlerByType.Load(meta.Type)

	if rawhandlers != nil {
		handlers := rawhandlers.([]cellnet.EventCallback)

		handlers = append(handlers, userCallback)
		self.handlerByType.Store(meta.Type, handlers)

	} else {
		self.handlerByType.Store(meta.Type, []cellnet.EventCallback{userCallback})
	}

}

var (
	GlobalDispatcher = new(MessageDispatcher)
)

func RegisterMessage(peer cellnet.Peer, msgName string, userHandler cellnet.EventCallback) {
	GlobalDispatcher.RegisterMessage(peer, msgName, userHandler)
}
