package socket

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

	if msgLogHook == nil || (msgLogHook != nil && msgLogHook(info)) {

		log.Debugf("#%s(%s) sid: %d %s size: %d | %s", info.Dir, info.PeerName, info.SessionID, info.Name, info.Size, info.Data)

	}

}

var msgLogHook func(*MessageLogInfo) bool

func SetMessageLogHook(hook func(*MessageLogInfo) bool) {
	msgLogHook = hook
}
