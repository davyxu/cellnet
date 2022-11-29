package cellrouter

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/x/container"
)

type Context struct {
	*xcontainer.Mapper
	Event cellnet.Event
}

func (self *Context) Message() any {
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

func (self *Context) MessageId() int {
	if self.Event == nil {
		return 0
	}

	return self.Event.MessageId()
}

func (self *Context) Session() cellnet.Session {
	if self.Event == nil {
		return nil
	}

	return self.Event.Session()
}

func (self *Context) Reset() {
	self.Event = nil
	self.Mapper = new(xcontainer.Mapper)
}

func (self *Context) Reply(msg any) {

	// 无法回复的来源, 多半是服务器间的来源
	if self.Event == nil {
		panic("can not reply to source due to nil sevent")
	}

	if replyEv, ok := self.Event.(interface {
		Reply(msg any)
	}); ok {
		replyEv.Reply(msg)
	} else {
		self.Session().Send(msg)
	}
}
