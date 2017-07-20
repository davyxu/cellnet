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
		ev.FromMessage(&coredef.SessionClosed{Result: r})
	case cellnet.Event_AcceptFailed:
		ev.FromMessage(&coredef.SessionAcceptFailed{Result: r})
	case cellnet.Event_ConnectFailed:
		ev.FromMessage(&coredef.SessionConnectFailed{Result: r})
	case cellnet.Event_Accepted:
		ev.MsgID = metaSessionAccepted.ID
	case cellnet.Event_Connected:
		ev.MsgID = metaSessionConnected.ID
	default:
		panic("unknown system error")
	}

	cellnet.MsgLog(ev)

	chain.Call(ev)
}
