package socket

import "github.com/davyxu/cellnet"

type MsgLogHandler struct {
	cellnet.BaseEventHandler
}

func (self *MsgLogHandler) Call(ev *cellnet.SessionEvent) (err error) {

	// 找到消息需要屏蔽
	if _, ok := msgMetaByID[ev.MsgID]; !ok {

		if msgLogHook == nil || (msgLogHook != nil && msgLogHook(ev)) {

			log.Debugf("#%s(%s) sid: %d %s size: %d | %s", ev.DirString(), ev.PeerName(), ev.SessionID(), ev.MsgName(), ev.MsgSize(), ev.MsgString())

		}
	}

	return self.CallNext(ev)
}

func NewMsgLogHandler() cellnet.EventHandler {

	return &MsgLogHandler{}

}

// 是否启用消息日志
var EnableMessageLog bool = true

var msgLogHook func(*cellnet.SessionEvent) bool
var msgMetaByID = make(map[uint32]*cellnet.MessageMeta)

func HookMessageLog(hook func(*cellnet.SessionEvent) bool) {
	msgLogHook = hook
}

func BlockMessageLog(msgName string) {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("msg log block not found: %s", msgName)
		return
	}

	msgMetaByID[meta.ID] = meta

}
