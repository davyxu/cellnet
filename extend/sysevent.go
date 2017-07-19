package extend

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/binary/coredef"
)

var (
	metaSessionConnected = cellnet.MessageMetaByName("coredef.SessionConnected")
	metaSessionAccepted  = cellnet.MessageMetaByName("coredef.SessionAccepted")
)

func PostSystemEvent(ses cellnet.Session, t cellnet.EventType, chain cellnet.HandlerChainList, r cellnet.Result) {

	ev := cellnet.NewEvent(t, ses)

	// 直接放在这里, decoder里遇到系统事件不会进行decode操作
	switch t {
	case cellnet.Event_Closed:
		ev.Msg = &coredef.SessionClosed{Result: r}
		ev.FromMessage(ev.Msg)
	case cellnet.Event_AcceptFailed:
		ev.Msg = &coredef.SessionAcceptFailed{Result: r}
		ev.FromMessage(ev.Msg)
	case cellnet.Event_ConnectFailed:
		ev.Msg = &coredef.SessionConnectFailed{Result: r}
		ev.FromMessage(ev.Msg)
	case cellnet.Event_Accepted:
		ev.FromMeta(metaSessionAccepted)
	case cellnet.Event_Connected:
		ev.FromMeta(metaSessionConnected)
	default:
		panic("unknown system error")
	}

	cellnet.MsgLog(ev)

	chain.Call(ev)
}
