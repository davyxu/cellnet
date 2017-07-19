package cellnet

type QueuePostHandler struct {
	q     EventQueue
	hlist []EventHandler
}

func (self *QueuePostHandler) Call(ev *Event) {

	self.q.Post(func() {
		HandlerChainCall(self.hlist, ev)
	})

}

func NewQueuePostHandler(q EventQueue, hlist ...EventHandler) EventHandler {

	return &QueuePostHandler{
		q:     q,
		hlist: hlist,
	}
}
