package cellevent

import "github.com/davyxu/cellnet"

// 回复对方消息,rpc
type Replier interface {
	Reply(msg interface{})
}

var (
	InternalDecodeHandler func(ev cellnet.Event) (msg interface{})
)
