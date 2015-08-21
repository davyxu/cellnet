package cellnet

type cell struct {
	id      CellID
	mailbox chan interface{}
}

// 投递到邮箱
func (self *cell) post(data interface{}) {

	// TODO 设置超时, 作为检查压力的一种方法
	self.mailbox <- data
}

func (self *cell) end() {
	self.mailbox <- endSignal{}
}

// 从邮箱取信
func (self *cell) fetch() (interface{}, bool) {

	data := <-self.mailbox

	if _, ok := data.(endSignal); ok {
		return nil, false
	}

	return data, true
}

type endSignal struct {
}
