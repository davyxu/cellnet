package socket

import (
	"reflect"

	"github.com/davyxu/cellnet"
)

func callSystemEvent(ses cellnet.Session, e cellnet.EventType, msg interface{}, h cellnet.EventHandler) {

	ev := cellnet.NewSessionEvent(e, nil)

	castToSystemEvent(ev, e, msg)

	cellnet.HandlerChainCall(h, ev)
}

func castToSystemEvent(ev *cellnet.SessionEvent, e cellnet.EventType, msg interface{}) {

	ev.Type = e

	meta := cellnet.MessageMetaByName(cellnet.MessageFullName(reflect.TypeOf(msg)))
	if meta != nil {
		ev.MsgID = meta.ID
	}

	// 直接放在这里, decoder里遇到系统事件不会进行decode操作
	ev.Msg = msg

}

func callSystemEventByMeta(ses cellnet.Session, e cellnet.EventType, meta *cellnet.MessageMeta, h cellnet.EventHandler) {

	ev := cellnet.NewSessionEvent(e, ses)

	ev.FromMeta(meta)

	cellnet.HandlerChainCall(h, ev)
}
