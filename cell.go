package cellnet

type cell struct {
	id      CellID
	target  CellID
	mailbox chan interface{}
}

func (self *cell) postMail(data interface{}) {

	self.mailbox <- data
}
