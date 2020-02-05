package msglog

import (
	"errors"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	whiteListByMsgID    sync.Map
	blackListByMsgID    sync.Map
	currMsgLogMode      = "black"
	currMsgLogModeGuard sync.RWMutex
)

// 设置当前的消息日志处理模式
// black: 黑名单模式, 黑名单中的消息不会显示, 其他均会显示
// white: 白名单模式, 只显示白名单中的消息, 其他不会显示
// mute: 屏蔽所有消息日志
// all: 显示所有消息日志
func SetCurrMsgLogMode(mode string) {
	currMsgLogModeGuard.Lock()
	currMsgLogMode = mode
	currMsgLogModeGuard.Unlock()
}

// 获取当前的消息日志处理模式
func GetCurrMsgLogMode() string {
	currMsgLogModeGuard.RLock()
	defer currMsgLogModeGuard.RUnlock()
	return currMsgLogMode
}

// 指定某个消息的处理规则, 消息格式: packageName.MsgName
// black: 黑名单模式, 黑名单中的消息不会显示, 其他均会显示
// white: 白名单模式, 只显示白名单中的消息, 其他不会显示
// none: 将此消息从白名单和黑名单中移除
func SetMsgLogRule(name string, rule string) error {

	meta := cellnet.MessageMetaByFullName(name)
	if meta == nil {
		return errors.New("msg not found")
	}

	switch rule {
	case "black":
		blackListByMsgID.Store(int(meta.ID), meta)
	case "white":
		whiteListByMsgID.Store(int(meta.ID), meta)
	case "none": // 从规则中移除
		blackListByMsgID.Delete(int(meta.ID))
		whiteListByMsgID.Delete(int(meta.ID))
	}

	return nil
}

// 能否显示消息日志
func IsMsgLogValid(msgid int) bool {
	switch GetCurrMsgLogMode() {
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

// 遍历消息规则
// black: 黑名单中的消息
// white: 白名单中的消息
func VisitMsgLogRule(mode string, callback func(*cellnet.MessageMeta) bool) {

	switch mode {
	case "black":
		blackListByMsgID.Range(func(key, value interface{}) bool {
			meta := value.(*cellnet.MessageMeta)

			return callback(meta)
		})
	case "white":
		whiteListByMsgID.Range(func(key, value interface{}) bool {
			meta := value.(*cellnet.MessageMeta)

			return callback(meta)
		})
	}

}
