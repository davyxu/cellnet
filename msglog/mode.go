package cellmsglog

import (
	"sync"
)

var (
	currMsgLogMode      = "black"
	currMsgLogModeGuard sync.RWMutex
)

// 设置当前的消息日志处理模式
// black: 黑名单模式, 黑名单中的消息不会显示, 其他均会显示
// white: 白名单模式, 只显示白名单中的消息, 其他不会显示
// mute: 屏蔽所有消息日志
// all: 显示所有消息日志
func SetMode(mode string) {
	currMsgLogModeGuard.Lock()
	currMsgLogMode = mode
	currMsgLogModeGuard.Unlock()
}

// 获取当前的消息日志处理模式
func Mode() string {
	currMsgLogModeGuard.RLock()
	defer currMsgLogModeGuard.RUnlock()
	return currMsgLogMode
}

// 能否显示消息日志
func IsMsgVisible(msgid int) bool {
	switch Mode() {
	case "black": // 黑名单里不显示
		if _, ok := blackListByMsgID.Load(msgid); ok {
			return false
		} else {
			return true
		}
	case "white": // 只有在白名单里才显示
		if _, ok := whiteListByMsgID.Load(msgid); ok {
			return true
		} else {
			return false
		}
	case "mute":
		return false
	case "all":
		return true
	default:
		return false
	}
}
