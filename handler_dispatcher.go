package cellnet

import (
	"bytes"
	"fmt"
	"sync"
)

type EventDispatcher interface {
	EventHandler

	// 注册事件回调
	AddHandler(id int, h []EventHandler) int

	// 根据AddHandler返回的int, 重新获得添加的Handler列表
	GetHandlerByIndex(id, index int) []EventHandler

	// 移除Handler
	RemoveHandler(id, index int)

	// 清除所有回调
	Clear()
}

type multiHandlerKey struct {
	id    int
	index int
}

func (self *multiHandlerKey) String() string {
	return fmt.Sprintf("id: %d index: %d", self.id, self.index)
}

type DispatcherHandler struct {
	handlerByKey      map[multiHandlerKey][]EventHandler
	handlerByKeyGuard sync.RWMutex
}

func (self *DispatcherHandler) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("\nDispatcher: %p\n", self))

	for key, hlist := range self.handlerByKey {

		buffer.WriteString(key.String())

		for _, h := range hlist {
			buffer.WriteString(fmt.Sprintf(" %v(%p) ", HandlerName(h), h))
		}

		buffer.WriteString("\n")

	}

	return buffer.String()
}

// 将连续调用的handler视为一个连接体 id -> Handler1 -> Handler1.1 -> Handler2 -> Handler 2.1

// 返回添加id对应的index, 删除需要两者
func (self *DispatcherHandler) AddHandler(id int, hlist []EventHandler) int {

	self.handlerByKeyGuard.Lock()

	key := self.findFreeIndex(id)
	self.handlerByKey[key] = hlist

	self.handlerByKeyGuard.Unlock()

	return key.index
}

func (self *DispatcherHandler) findFreeIndex(id int) multiHandlerKey {

	key := multiHandlerKey{id, 0}

	for index := 0; ; index++ {

		key.index = index

		if v, ok := self.handlerByKey[key]; !ok || v == nil {
			return key
		}
	}

}

func (self *DispatcherHandler) Call(ev *Event) {

	key := multiHandlerKey{int(ev.MsgID), 0}

	for index := 0; ; index++ {

		// 拷贝一份, 放置派发时, 内部被修改
		copyed := ev.Clone()

		key.index = index

		self.handlerByKeyGuard.RLock()
		hlist, ok := self.handlerByKey[key]
		self.handlerByKeyGuard.RUnlock()

		if ok {

			HandlerChainCall(hlist, copyed)
		} else {
			break
		}
	}

}

// 移除handler
func (self *DispatcherHandler) RemoveHandler(id, index int) {

	self.handlerByKeyGuard.Lock()
	self.handlerByKey[multiHandlerKey{id, index}] = nil
	self.handlerByKeyGuard.Unlock()
}

func (self *DispatcherHandler) Clear() {

	self.handlerByKeyGuard.Lock()
	self.handlerByKey = make(map[multiHandlerKey][]EventHandler)
	self.handlerByKeyGuard.Unlock()
}

// index 根据注册顺序, 从0~n
func (self *DispatcherHandler) GetHandlerByIndex(id, index int) []EventHandler {

	self.handlerByKeyGuard.RLock()
	hlist, ok := self.handlerByKey[multiHandlerKey{id, index}]
	self.handlerByKeyGuard.RUnlock()

	if ok {
		return hlist
	}

	return nil
}

func NewDispatcherHandler() *DispatcherHandler {
	self := &DispatcherHandler{}

	self.Clear()

	return self
}
