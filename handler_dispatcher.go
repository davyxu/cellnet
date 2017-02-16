package cellnet

type EventDispatcher interface {

	// 注册事件回调
	AddHandler(id int, h EventHandler) *HandlerContext

	RemoveHandler(id int)

	Call(*SessionEvent) error

	// 清除所有回调
	Clear()

	Count() int

	CountByID(id int) int

	VisitCallback(callback func(int, *HandlerContext) VisitOperation)
}

type HandlerContext struct {
	ID      int
	Handler EventHandler

	Tag interface{}
}

type DispatcherHandler struct {
	BaseEventHandler
	handlerByID map[int][]*HandlerContext
}

func (self *DispatcherHandler) AddHandler(id int, h EventHandler) *HandlerContext {

	// 事件
	ctxList, ok := self.handlerByID[id]

	if !ok {
		ctxList = make([]*HandlerContext, 0)
	}

	newCtx := &HandlerContext{
		ID:      id,
		Handler: h,
	}

	ctxList = append(ctxList, newCtx)

	self.handlerByID[int(id)] = ctxList

	return newCtx
}

func (self *DispatcherHandler) RemoveHandler(id int) {

	delete(self.handlerByID, id)
}

func (self *DispatcherHandler) Call(ev *SessionEvent) error {

	if ctxList, ok := self.handlerByID[int(ev.MsgID)]; ok {

		for _, ctx := range ctxList {

			if err := HandlerCallNext(ctx.Handler, ev); err != nil {
				return err
			}
		}

	}

	return nil
}

type VisitOperation int

const (
	VisitOperation_Continue = iota // 循环下一个
	VisitOperation_Remove          // 删除当前元素
	VisitOperation_Exit            // 退出循环
)

func (self *DispatcherHandler) VisitCallback(callback func(int, *HandlerContext) VisitOperation) {

	var needDelete []int

	for id, ctxList := range self.handlerByID {

		var needRefresh bool

		var index = 0
		for {

			if index >= len(ctxList) {
				break
			}

			ctx := ctxList[index]

			op := callback(id, ctx)

			switch op {
			case VisitOperation_Exit:
				goto endloop
			case VisitOperation_Remove:

				if len(ctxList) == 1 {
					needDelete = append(needDelete, id)
				}

				ctxList = append(ctxList[:index], ctxList[index+1:]...)

				needRefresh = true
			case VisitOperation_Continue:
				index++
			}

		}

		if needRefresh {
			self.handlerByID[id] = ctxList
		}

	}

endloop:

	if len(needDelete) > 0 {
		for _, id := range needDelete {
			delete(self.handlerByID, id)
		}
	}

}

func (self *DispatcherHandler) Clear() {

	self.handlerByID = make(map[int][]*HandlerContext)
}

func (self *DispatcherHandler) Exists(id int) bool {

	_, ok := self.handlerByID[id]

	return ok
}

func (self *DispatcherHandler) Count() int {
	return len(self.handlerByID)
}

func (self *DispatcherHandler) CountByID(id int) int {

	if v, ok := self.handlerByID[id]; ok {
		return len(v)
	}

	return 0
}

func NewDispatcherHandler() *DispatcherHandler {
	return &DispatcherHandler{
		handlerByID: make(map[int][]*HandlerContext),
	}
}
