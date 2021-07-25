package cellrouter

import (
	"github.com/davyxu/cellnet"
	xframe "github.com/davyxu/x/frame"
	"math"
)

type HandlerFunc func(ctx *Context)

type Context struct {
	*xframe.PropertySet
	Event cellnet.Event

	index int8

	handlers []HandlerFunc
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
		panic("Require 'ReplyEvent' to reply event")
	}
}

func (self *Context) Reset() {
	self.handlers = nil
	self.index = -1
	self.Event = nil
	self.PropertySet = new(xframe.PropertySet)
}

const abortIndex int8 = math.MaxInt8 / 2

func (self *Context) Next() {
	self.index++
	for self.index < int8(len(self.handlers)) {
		self.handlers[self.index](self)
		self.index++
	}
}

func (self *Context) Abort() {
	self.index = abortIndex
}
func (self *Context) IsAborted() bool {
	return self.index == abortIndex
}
