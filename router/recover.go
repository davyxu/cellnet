package cellrouter

import (
	cellmeta "github.com/davyxu/cellnet/meta"
	xlog "github.com/davyxu/x/logger"
	xruntime "github.com/davyxu/x/runtime"
)

// 用法: Global.Use(Recover())
func Recover() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {

				reqName := cellmeta.MessageToName(ctx.Message())
				reqBody := cellmeta.MessageToString(ctx.Message())

				xlog.Errorf("Panic recovery | %v | stack: %s |> %s | %s", err, xruntime.StackToString(5), reqName, reqBody)
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
