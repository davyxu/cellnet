package cellnet

import (
	"errors"
	"sync"
)

type MsgLogHandler struct {
}

func (self *MsgLogHandler) Call(ev *Event) {

	MsgLog(ev)
}

var defaultmsgLogHandler = new(MsgLogHandler)

func StaticMsgLogHandler() EventHandler {
	return defaultmsgLogHandler
}

// Msg
// Data, MsgID

func MsgLog(ev *Event) {

	if !log.IsDebugEnabled() {
		return
	}

	ev.Parse()

	if IsBlockedMessageByID(ev.MsgID) {
		return
	}

	// 需要在收到消息, 不经过decoder时, 就要打印出来, 所以手动解开消息, 有少许耗费

	log.Debugf("#%s(%s) sid: %d %s(%d) size: %d | %s", ev.Type.String(), ev.PeerName(), ev.SessionID(), ev.MsgName(), ev.MsgID, ev.MsgSize(), ev.MsgString())

}

var (
	msgMetaByID      = map[uint32]*MessageMeta{}
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
	meta := MessageMetaByName(msgName)

	if meta == nil {
		return ErrMessageNotFound
	}

	msgMetaByIDGuard.Lock()
	msgMetaByID[meta.ID] = meta
	msgMetaByIDGuard.Unlock()

	return nil
}
