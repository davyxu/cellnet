package cellmsglog

import (
	"errors"
	cellmeta "github.com/davyxu/cellnet/meta"
	"sync"
)

var (
	whiteListByMsgID sync.Map
	blackListByMsgID sync.Map
)

// 指定某个消息的处理规则, 消息格式: packageName.MsgName
// black: 黑名单模式, 黑名单中的消息不会显示, 其他均会显示
// white: 白名单模式, 只显示白名单中的消息, 其他不会显示
// none: 将此消息从白名单和黑名单中移除
func SetRule(name string, rule string) error {

	meta := cellmeta.MetaByFullName(name)
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

// 遍历消息规则
// black: 黑名单中的消息
// white: 白名单中的消息
func VisitRule(mode string, callback func(*cellmeta.Meta) bool) {

	switch mode {
	case "black":
		blackListByMsgID.Range(func(key, value any) bool {
			meta := value.(*cellmeta.Meta)

			return callback(meta)
		})
	case "white":
		whiteListByMsgID.Range(func(key, value any) bool {
			meta := value.(*cellmeta.Meta)

			return callback(meta)
		})
	}

}
