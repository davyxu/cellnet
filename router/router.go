package cellrouter

import (
	cellmeta "github.com/davyxu/cellnet/meta"
	xframe "github.com/davyxu/x/frame"
)

// gin风格的路由
type Router struct {
	node            xframe.PropertySet
	handlers        []HandlerFunc // 全局handler
	defaultHandlers []HandlerFunc // 默认无处理handler
}

func (self *Router) Use(handlers ...HandlerFunc) {
	self.handlers = append(self.handlers, handlers...)
}

func (self *Router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(self.handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, self.handlers)
	copy(mergedHandlers[len(self.handlers):], handlers)
	return mergedHandlers
}

func (self *Router) Handle(msgTypeObj interface{}, handlers ...HandlerFunc) {

	meta := cellmeta.MetaByMsg(msgTypeObj)
	if meta == nil {
		panic("msg not register meta")
	}

	self.node.Set(meta.ID, self.combineHandlers(handlers))
}

func (self *Router) HandleDefault(handlers ...HandlerFunc) {
	self.defaultHandlers = self.combineHandlers(handlers)
}

func (self *Router) Invoke(ctx *Context) {
	if raw, ok := self.node.Get(ctx.MessageID()); ok {
		ctx.handlers = raw.([]HandlerFunc)
		ctx.Next()
	} else if len(self.defaultHandlers) > 0 {
		ctx.handlers = self.defaultHandlers
		ctx.Next()
	}
}
