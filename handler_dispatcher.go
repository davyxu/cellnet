package cellnet

type EventDispatcher interface {

	// 注册事件回调
	AddHandler(id int, h EventHandler)

	GetHandlerByID(id int) EventHandler

	RemoveHandler(id int)

	Call(*SessionEvent)

	// 清除所有回调
	Clear()

	Count() int

	CountByID(id int) int
}

type DispatcherHandler struct {
	BaseEventHandler
	handlerByID map[int]EventHandler

	headerHandler map[EventHandler]bool
}

// 将连续调用的handler视为一个连接体 id -> Handler1 -> Handler1.1 -> Handler2 -> Handler 2.1

func (self *DispatcherHandler) GetHandlerByID(id int) EventHandler {
	if exists, ok := self.handlerByID[int(id)]; ok {
		return exists
	}

	return nil
}

func (self *DispatcherHandler) AddHandler(id int, h EventHandler) {

	// 回调不允许在多个id被注册
	if _, ok := self.headerHandler[h]; ok {
		panic("Duplicate header handler")
	}

	self.headerHandler[h] = true

	if exists, ok := self.handlerByID[int(id)]; ok {

		// 找到尾巴
		head := findTail(exists)

		// 连上去
		if head != nil {
			head.SetNext(h)
			return
		}

	}

	self.handlerByID[int(id)] = h

}

func findTail(origin EventHandler) EventHandler {

	h := origin

	for h != nil {

		if h.Next() == nil {
			return h
		}
	}

	return nil

}

func (self *DispatcherHandler) RemoveHandler(id int) {

	delete(self.handlerByID, id)
}

func (self *DispatcherHandler) Next() EventHandler {
	return nil
}

func (self *DispatcherHandler) Call(ev *SessionEvent) {

	if h, ok := self.handlerByID[int(ev.MsgID)]; ok {

		HandlerChainCall(h, ev)
	}

}

func (self *DispatcherHandler) Clear() {

	self.handlerByID = make(map[int]EventHandler)
}

func (self *DispatcherHandler) Exists(id int) bool {

	_, ok := self.handlerByID[id]

	return ok
}

func (self *DispatcherHandler) Count() int {
	return len(self.handlerByID)
}

func (self *DispatcherHandler) CountByID(id int) int {

	if _, ok := self.handlerByID[id]; ok {
		return 1
	}

	return 0
}

func NewDispatcherHandler() *DispatcherHandler {
	return &DispatcherHandler{
		handlerByID:   make(map[int]EventHandler),
		headerHandler: make(map[EventHandler]bool),
	}
}
