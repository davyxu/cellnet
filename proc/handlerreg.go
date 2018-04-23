package proc

import (
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

// 消息派发器，可选件，兼容v3以前的注册及派发消息方式，在没有代码生成框架及工具时是较方便的接收处理接口
type MessageDispatcher struct {
	handlerByType      map[reflect.Type][]cellnet.EventCallback
	handlerByTypeGuard sync.RWMutex
}

func (self *MessageDispatcher) OnEvent(ev cellnet.Event) {

	msgType := reflect.TypeOf(ev.Message())

	self.handlerByTypeGuard.RLock()
	handlers, ok := self.handlerByType[msgType.Elem()]
	self.handlerByTypeGuard.RUnlock()

	if ok {

		for _, callback := range handlers {

			callback(ev)
		}

	}
}

func (self *MessageDispatcher) RegisterMessage(msgName string, userCallback cellnet.EventCallback) {
	meta := cellnet.MessageMetaByFullName(msgName)
	if meta == nil {
		panic("message not found:" + msgName)
	}

	self.handlerByTypeGuard.Lock()
	handlers, _ := self.handlerByType[meta.Type]
	handlers = append(handlers, userCallback)
	self.handlerByType[meta.Type] = handlers
	self.handlerByTypeGuard.Unlock()
}

func NewMessageDispatcher(peer cellnet.Peer, processorName string) *MessageDispatcher {

	self := &MessageDispatcher{
		handlerByType: make(map[reflect.Type][]cellnet.EventCallback),
	}

	BindProcessorHandler(peer, processorName, self.OnEvent)

	return self
}
