package cellnet

type EventPipe interface {
	AddQueue() EventQueue

	Start()

	Stop(int)

	Wait() int
}

type evPipe struct {
	qarray []*evQueue

	arrayLock  bool
	exitSignal chan int
}

func (self *evPipe) AddQueue() EventQueue {

	if self.arrayLock {
		panic("Pipe already start, can not addqueue any more")
	}

	q := newEventQueue()

	self.qarray = append(self.qarray, q)

	return q
}

type combinedEvent struct {
	q *evQueue
	e interface{}
}

func (self *evPipe) Start() {

	// 开始后, 不能修改数组
	self.arrayLock = true

	go func() {

		combinedChannel := make(chan *combinedEvent)

		for _, q := range self.qarray {
			go func(q *evQueue) {
				for v := range q.queue {
					combinedChannel <- &combinedEvent{q: q, e: v}
				}
			}(q)
		}

		for v := range combinedChannel {
			v.q.CallData(v.e)
		}

	}()

}

func (self *evPipe) Stop(result int) {
	self.exitSignal <- result
}

func (self *evPipe) Wait() int {
	return <-self.exitSignal
}

func NewEventPipe() EventPipe {
	return &evPipe{
		qarray:     make([]*evQueue, 0),
		exitSignal: make(chan int),
	}
}
