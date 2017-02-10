package socket

import "github.com/davyxu/cellnet"

type DispatcherHandler struct {
	handlerByID map[uint32][]cellnet.Handler
}

func (self *DispatcherHandler) Add(id uint32, h cellnet.Handler) {

	// 事件
	handlers, ok := self.handlerByID[id]

	if !ok {
		handlers = make([]cellnet.Handler, 0)

	}

	handlers = append(handlers, h)

	self.handlerByID[id] = handlers
}

func (self *DispatcherHandler) Call(evid int, data interface{}) error {

	ev := data.(*SessionEvent)

	if handlers, ok := self.handlerByID[ev.MsgID]; ok {

		for _, h := range handlers {
			if err := h.Call(evid, data); err != nil {
				return err
			}
		}

	}

	return nil
}

func NewDispatcherHandler() *DispatcherHandler {
	return &DispatcherHandler{
		handlerByID: make(map[uint32][]cellnet.Handler),
	}
}
