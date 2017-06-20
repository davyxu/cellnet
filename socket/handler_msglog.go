package socket

import (
	"fmt"

	"errors"
	"github.com/davyxu/cellnet"
	"sync"
)

type MsgLogHandler struct {
}

func (self *MsgLogHandler) Call(ev *cellnet.Event) {

	MsgLog(ev)
}

var defaultmsgLogHandler = new(MsgLogHandler)

func StaticMsgLogHandler() cellnet.EventHandler {
	return defaultmsgLogHandler
}

// Msg
// Data, MsgID

func MsgLog(ev *cellnet.Event) {

	ev.Parse()

	if IsBlockedMessageByID(ev.MsgID) {
		return
	}

	// 需要在收到消息, 不经过decoder时, 就要打印出来, 所以手动解开消息, 有少许耗费

	log.Debugf("#%s(%s) sid: %d %s(%d) size: %d | %s", dirString(ev), ev.PeerName(), ev.SessionID(), ev.MsgName(), ev.MsgID, ev.MsgSize(), ev.MsgString())

}

func dirString(ev *cellnet.Event) string {

	switch ev.Type {
	case cellnet.Event_Recv:
		return "recv"
	case cellnet.Event_Post:
		return "post"
	case cellnet.Event_Send:
		return "send"
	case cellnet.Event_Connected:
		return "connected"
	case cellnet.Event_ConnectFailed:
		return "connectfailed"
	case cellnet.Event_Accepted:
		return "accepted"
	case cellnet.Event_AcceptFailed:
		return "acceptefailed"
	case cellnet.Event_Closed:
		return "closed"
	}

	return fmt.Sprintf("unknown(%d)", ev.Type)
}

var (

	// 是否启用消息日志
	EnableMessageLog bool = true

	msgMetaByID      = map[uint32]*cellnet.MessageMeta{}
	msgMetaByIDGuard sync.RWMutex
)

func IsBlockedMessageByID(msgid uint32) bool {
	msgMetaByIDGuard.RLock()
	defer msgMetaByIDGuard.RUnlock()

	if _, ok := msgMetaByID[msgid]; ok {
		return true
	}

	return false
}

var (
	ErrMessageNotFound = errors.New("msg not exists")
)

func BlockMessageLog(msgName string) error {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		return ErrMessageNotFound
	}

	msgMetaByIDGuard.Lock()
	msgMetaByID[meta.ID] = meta
	msgMetaByIDGuard.Unlock()

	return nil
}
