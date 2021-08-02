package cellrouter

import (
	"github.com/davyxu/cellnet"
	xframe "github.com/davyxu/x/frame"
)

type Context struct {
	*xframe.Mapper
	Event cellnet.Event
}

func (self *Context) Message() interface{} {
	if self.Event == nil {
		return nil
	}

	return self.Event.Message()
}

func (self *Context) MessageData() []byte {
	if self.Event == nil {
		return nil
	}

	return self.Event.MessageData()
}

func (self *Context) MessageID() int {
	if self.Event == nil {
		return 0
	}

	return self.Event.MessageID()
}

func (self *Context) Session() cellnet.Session {
	if self.Event == nil {
		return nil
	}

	return self.Event.Session()
}

func (self *Context) Reset() {
	self.Event = nil
	self.Mapper = new(xframe.Mapper)
}
