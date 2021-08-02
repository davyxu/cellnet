package cellrouter

import (
	cellmeta "github.com/davyxu/cellnet/meta"
	xframe "github.com/davyxu/x/frame"
	xlog "github.com/davyxu/x/logger"
	xos "github.com/davyxu/x/os"
	xruntime "github.com/davyxu/x/runtime"
)

type HandlerFunc func(ctx *Context)

type Router struct {
	mapper   xframe.Mapper
	handlers []interface{} // 全局handler
	Recover  bool
}

type HandlerKey struct {
	ID   int
	Kind string
}

func (self *Router) Handle(obj interface{}, kind string, handler interface{}) {

	meta := cellmeta.MetaByMsg(obj)
	if meta == nil {
		panic("msg not register meta")
	}
	self.mapper.Set(HandlerKey{ID: meta.ID, Kind: kind}, handler)
}

func (self *Router) Invoke(ctx *Context, kind string, customInvoker func(raw interface{})) {

	if self.Recover {
		defer xos.Recover(func(raw interface{}) {
			reqName := cellmeta.MessageToName(ctx.Message())
			reqBody := cellmeta.MessageToString(ctx.Message())

			xlog.Errorf("Panic recovery | %v | stack: %s |> %s | %s", raw, xruntime.StackToString(5), reqName, reqBody)
		})
	}

	if raw, ok := self.mapper.Get(HandlerKey{ID: ctx.MessageID(), Kind: kind}); ok {
		if entry, ok := raw.(HandlerFunc); ok {

			entry(ctx)
		} else if customInvoker != nil {
			customInvoker(raw)
		}
	}
}
