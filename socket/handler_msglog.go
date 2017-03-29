package socket

import (
	"fmt"

	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

type MsgLogHandler struct {
	cellnet.BaseEventHandler
}

func dirString(ev *cellnet.SessionEvent) string {

	switch ev.Type {
	case cellnet.SessionEvent_Recv:
		return "recv"
	case cellnet.SessionEvent_Post:
		return "post"
	case cellnet.SessionEvent_Send:
		return "send"
	case cellnet.SessionEvent_Connected:
		return "connected"
	case cellnet.SessionEvent_ConnectFailed:
		return "connectfailed"
	case cellnet.SessionEvent_Accepted:
		return "accepted"
	case cellnet.SessionEvent_AcceptFailed:
		return "acceptefailed"
	case cellnet.SessionEvent_Closed:
		return "closed"
	}

	return fmt.Sprintf("unknown(%d)", ev.Type)
}

func (self *MsgLogHandler) Call(ev *cellnet.SessionEvent) {

	MsgLog(ev)
}

// Msg
// Data, MsgID

func MsgLog(ev *cellnet.SessionEvent) {

	// send和post时
	if ev.MsgID == 0 && ev.Msg != nil {
		meta := cellnet.MessageMetaByName(cellnet.MessageFullName(reflect.TypeOf(ev.Msg)))
		if meta != nil {
			ev.MsgID = meta.ID
		}
	} else {
		// 接收时

		ev.Msg, _ = cellnet.DecodeMessage(ev.MsgID, ev.Data)
	}

	if IsBlockedMessageByID(ev.MsgID) {
		return
	}

	if msgLogHook == nil || (msgLogHook != nil && msgLogHook(ev)) {

		// 需要在收到消息, 不经过decoder时, 就要打印出来, 所以手动解开消息, 有少许耗费

		log.Debugf("#%s(%s) sid: %d %s size: %d | %s", dirString(ev), ev.PeerName(), ev.SessionID(), ev.MsgName(), ev.MsgSize(), ev.MsgString())

	}

}

func NewMsgLogHandler() cellnet.EventHandler {

	return &MsgLogHandler{}

}

var (

	// 是否启用消息日志
	EnableMessageLog bool = true

	msgLogHook       func(*cellnet.SessionEvent) bool
	msgMetaByID      = map[uint32]*cellnet.MessageMeta{}
	msgMetaByIDGuard sync.RWMutex
)

func HookMessageLog(hook func(*cellnet.SessionEvent) bool) {
	msgLogHook = hook
}

func IsBlockedMessageByID(msgid uint32) bool {
	msgMetaByIDGuard.RLock()
	defer msgMetaByIDGuard.RUnlock()

	if _, ok := msgMetaByID[msgid]; ok {
		return true
	}

	return false
}

func BlockMessageLog(msgName string) {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("msg log block not found: %s", msgName)
		return
	}

	msgMetaByIDGuard.Lock()
	msgMetaByID[meta.ID] = meta
	msgMetaByIDGuard.Unlock()

}
