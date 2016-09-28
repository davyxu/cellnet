package socket

import (
	"github.com/davyxu/cellnet"
)

type MessageLogInfo struct {
	Dir       string
	PeerName  string
	SessionID int64
	Name      string
	ID        uint32
	Size      int32
	Data      string
}

// 是否启用消息日志
var EnableMessageLog bool = true

func msgLog(info *MessageLogInfo) {

	// 找到消息需要屏蔽
	if _, ok := msgMetaByID[info.ID]; ok {
		return
	}

	if msgLogHook == nil || (msgLogHook != nil && msgLogHook(info)) {

		log.Debugf("#%s(%s) sid: %d %s size: %d | %s", info.Dir, info.PeerName, info.SessionID, info.Name, info.Size, info.Data)

	}

}

var msgLogHook func(*MessageLogInfo) bool

func SetMessageLogHook(hook func(*MessageLogInfo) bool) {
	msgLogHook = hook
}

var msgMetaByID = make(map[uint32]*cellnet.MessageMeta)

func BlockMessageLog(msgName string) {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		log.Errorf("msg log block not found: %s", msgName)
		return
	}

	msgMetaByID[meta.ID] = meta

}
