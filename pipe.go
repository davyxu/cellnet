package cellnet

type EventPipe interface {
	AddQueue() EventQueue

	Start()

	Stop(int)

	Wait() int
}

type evPipe struct {
	exitSignal chan int

	dataChan chan *dataTask
}

func (self *evPipe) AddQueue() EventQueue {

	q := newEventQueue()

	go func(q *evQueue) {

		for v := range q.queue {
			self.dataChan <- &dataTask{q: q, e: v}
		}

	}(q)

	return q
}

type dataTask struct {
	q *evQueue
	e interface{}
}

func (self *evPipe) Start() {

	go func() {

		for v := range self.dataChan {
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
		exitSignal: make(chan int),
		dataChan:   make(chan *dataTask),
	}
}
