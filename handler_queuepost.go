package cellnet

type QueuePostHandler struct {
	BaseEventHandler
	q EventQueue
}

func (self *QueuePostHandler) Next() EventHandler {
	return nil
}

func (self *QueuePostHandler) Call(ev *SessionEvent) {

	self.q.Post(func() {

		HandlerChainCall(self.next, ev)
	})

}

func NewQueuePostHandler(q EventQueue) EventHandler {
	return &QueuePostHandler{
		q: q,
	}
}
