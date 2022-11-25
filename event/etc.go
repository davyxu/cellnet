package cellevent

import "github.com/davyxu/cellnet"

var (
	InternalDecodeHandler func(ev cellnet.Event) (msg any)
)
