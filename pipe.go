package cellnet

import (
	"reflect"
)

type EvPipe struct {
	qarray []*EvQueue

	arrayLock bool
}

func (self *EvPipe) AddQueue() *EvQueue {

	if self.arrayLock {
		panic("Pipe already start, can not addqueue any more")
	}

	q := newEvQueue()

	self.qarray = append(self.qarray, q)

	return q
}

func (self *EvPipe) Start() {

	// 开始后, 不能修改数组
	self.arrayLock = true

	go func() {

		cases := make([]reflect.SelectCase, len(self.qarray))

		// 按照队列(peer)数量开始做case
		for i, q := range self.qarray {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(q.queue)}
		}

		for {

			if index, value, ok := reflect.Select(cases); ok {

				self.qarray[index].call(value.Interface())
			}

		}

	}()

}

func NewEvPipe() *EvPipe {
	return &EvPipe{
		qarray: make([]*EvQueue, 0),
	}
}
