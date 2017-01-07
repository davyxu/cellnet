package cellnet

type EventDispatcher interface {

	// 注册事件回调
	RegisterCallback(id uint32, f func(interface{}))

	// 设置事件截获钩子, 在CallData中调用钩子
	InjectData(func(interface{}) bool)

	// 直接调用消费者端的handler
	CallData(data interface{})
}

type evDispatcher struct {
	// 保证注册发生在初始化, 读取发生在之后可以不用锁
	handlerByMsgPeer map[uint32][]func(interface{})

	inject func(interface{}) bool
}

// 注册事件回调
func (self *evDispatcher) RegisterCallback(id uint32, f func(interface{})) {

	// 事件
	em, ok := self.handlerByMsgPeer[id]

	// 新建
	if !ok {

		em = make([]func(interface{}), 0)

	}

	em = append(em, f)

	self.handlerByMsgPeer[id] = em
}

// 注入回调, 返回false时表示不再投递
func (self *evDispatcher) InjectData(f func(interface{}) bool) {

	self.inject = f
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

		if carr, ok := self.handlerByMsgPeer[d.ContextID()]; ok {

			// 遍历所有的回调
			for _, c := range carr {

				c(data)
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
		handlerByMsgPeer: make(map[uint32][]func(interface{})),
	}

	return self

}
