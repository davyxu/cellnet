package cellevent

import "github.com/davyxu/cellnet"

// rpc, relay, 普通消息
type ReplyEvent interface {
	Reply(msg interface{})
}

var (
	InternalDecodeHandler func(ev cellnet.Event) (msg interface{})
)
