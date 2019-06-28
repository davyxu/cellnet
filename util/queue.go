package util

type Queue struct {
	list []interface{}
}

func (self *Queue) Enqueue(data interface{}) {

	self.list = append(self.list, data)
}

func (self *Queue) Count() int {
	return len(self.list)
}

func (self *Queue) Peek() interface{} {
	return self.list[0]
}

func (self *Queue) Dequeue() (ret interface{}) {

	if len(self.list) == 0 {
		return nil
	}

	ret = self.list[0]

	self.list = self.list[1:]

	return
}

func (self *Queue) Clear() {
	self.list = self.list[0:0]
}

func NewQueue(size int) *Queue {

	return &Queue{
		list: make([]interface{}, 0, size),
	}
}
