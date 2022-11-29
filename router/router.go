package cellrouter

import (
	cellmeta "github.com/davyxu/cellnet/meta"
	"github.com/davyxu/x/container"
	xos "github.com/davyxu/x/os"
	xruntime "github.com/davyxu/x/runtime"
	"github.com/davyxu/xlog"
)

type HandlerFunc func(ctx *Context)

type Router struct {
	mapper   xcontainer.Mapper
	handlers []any // 全局handler
	Recover  bool
}

type HandlerKey struct {
	ID   int
	Kind string
}

func (self *Router) Handle(obj any, kind string, handler any) {

	meta := cellmeta.MetaByMsg(obj)
	if meta == nil {
		panic("msg not register meta")
	}
	self.mapper.Set(HandlerKey{ID: meta.Id, Kind: kind}, handler)
}

func (self *Router) Invoke(ctx *Context, kind string, customInvoker func(raw any)) {

	if self.Recover {
		defer xos.Recover(func(raw any) {
			reqName := cellmeta.MessageToName(ctx.Message())
			reqBody := cellmeta.MessageToString(ctx.Message())

			xlog.Errorf("Panic recovery | %v | stack: %s |> %s | %s", raw, xruntime.StackToString(5), reqName, reqBody)
		})
	}

	if raw, ok := self.mapper.Get(HandlerKey{ID: ctx.MessageId(), Kind: kind}); ok {
		if entry, ok := raw.(HandlerFunc); ok {

			entry(ctx)
		} else if customInvoker != nil {
			customInvoker(raw)
		}
	}
}
