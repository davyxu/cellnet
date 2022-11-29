package cellcodec

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
)

func init() {
	cellevent.InternalDecodeHandler = func(ev cellnet.Event) (msg any) {
		msg, _, _ = Decode(ev.MessageId(), ev.MessageData())
		return
	}
}
