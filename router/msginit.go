package cellrouter

import (
	cellevent "github.com/davyxu/cellnet/event"
	"github.com/davyxu/cellnet/peer/tcp"
)

var (
	Global = new(Router)
)

// 注册消息处理具备
func Handle(msgTypeObj interface{}, handler HandlerFunc) {
	Global.Handle(msgTypeObj, "", handler)
}

// 对接cellnet.Peer.OnInbound
func InboundEntry(input *cellevent.RecvMsg) (output *cellevent.RecvMsg) {

	ctx := new(Context)
	ctx.Reset()
	ctx.Event = input

	RecvLogger(tcp.SessionID(input.Session()), ctx)

	Global.Invoke(ctx, "", nil)

	return input

}
