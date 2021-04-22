package cellcodec

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
)

func init() {
	cellevent.InternalDecodeHandler = func(ev cellnet.Event) (msg interface{}) {
		msg, _, _ = Decode(ev.MessageID(), ev.MessageData())
		return
	}
}
