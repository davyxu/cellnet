package cellnet

type EventDispatcher interface {

	// 注册事件回调
	AddCallback(id uint32, f func(interface{})) *CallbackContext

	RemoveCallback(id uint32)

	// 设置事件截获钩子, 在CallData中调用钩子
	InjectData(func(interface{}) bool)

	// 直接调用消费者端的handler
	CallData(data interface{})

	// 清除所有回调
	Clear()

	Count() int

	CountByID(id uint32) int

	VisitCallback(callback func(uint32, *CallbackContext) VisitOperation)
}

type CallbackContext struct {
	ID   uint32
	Func func(interface{})

	Tag interface{}
}

type evDispatcher struct {
	// 保证注册发生在初始化, 读取发生在之后可以不用锁
	handlerByMsgPeer map[uint32][]*CallbackContext

	inject func(interface{}) bool
}

// 注册事件回调
func (self *evDispatcher) AddCallback(id uint32, f func(interface{})) *CallbackContext {

	// 事件
	ctxList, ok := self.handlerByMsgPeer[id]

	if !ok {
		ctxList = make([]*CallbackContext, 0)

	}

	newCtx := &CallbackContext{
		ID:   id,
		Func: f,
	}

	ctxList = append(ctxList, newCtx)

	self.handlerByMsgPeer[id] = ctxList

	return newCtx
}

func (self *evDispatcher) RemoveCallback(id uint32) {

	delete(self.handlerByMsgPeer, id)
}

// 注入回调, 返回false时表示不再投递
func (self *evDispatcher) InjectData(f func(interface{}) bool) {

	self.inject = f
}

type VisitOperation int

const (
	VisitOperation_Continue = iota // 循环下一个
	VisitOperation_Remove          // 删除当前元素
	VisitOperation_Exit            // 退出循环
)

func (self *evDispatcher) VisitCallback(callback func(uint32, *CallbackContext) VisitOperation) {

	var needDelete []uint32

	for id, ctxList := range self.handlerByMsgPeer {

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
			self.handlerByMsgPeer[id] = ctxList
		}

	}

endloop:

	if len(needDelete) > 0 {
		for _, id := range needDelete {
			delete(self.handlerByMsgPeer, id)
		}
	}

}

func (self *evDispatcher) Clear() {

	self.handlerByMsgPeer = make(map[uint32][]*CallbackContext)
}

func (self *evDispatcher) Exists(id uint32) bool {

	_, ok := self.handlerByMsgPeer[id]

	return ok
}

func (self *evDispatcher) Count() int {
	return len(self.handlerByMsgPeer)
}

func (self *evDispatcher) CountByID(id uint32) int {

	if v, ok := self.handlerByMsgPeer[id]; ok {
		return len(v)
	}

	return 0
}

type contentIndexer2 interface {
	ContextID() uint32
}

// 通过数据接口调用
func (self *evDispatcher) CallData(data interface{}) {

	switch d := data.(type) {
	// ID索引的消息
	case contentIndexer2:

		if self == nil {
			log.Errorf("recv indexed event, but event dispatcher nil, id: %d", d.ContextID())
			return
		}

		// 先处理注入
		if self.inject != nil && !self.inject(data) {
			return
		}

		if ctxList, ok := self.handlerByMsgPeer[d.ContextID()]; ok {

			for _, ctx := range ctxList {
				ctx.Func(data)
			}

		}
	// 直接回调
	case func():
		d()
	default:
		log.Errorln("unknown queue data: ", data)
	}

}

func NewEventDispatcher() EventDispatcher {
	self := &evDispatcher{
		handlerByMsgPeer: make(map[uint32][]*CallbackContext),
	}

	return self

}
