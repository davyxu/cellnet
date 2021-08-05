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

func (self *Context) Reply(msg interface{}) {

	// 无法回复的来源, 多半是服务器间的来源
	if self.Event == nil {
		panic("can not reply to source due to nil sevent")
	}

	if replyEv, ok := self.Event.(interface {
		Reply(msg interface{})
	}); ok {
		replyEv.Reply(msg)
	} else {
		self.Session().Send(msg)
	}
}
