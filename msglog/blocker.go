package msglog

import (
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	blockedMsgByID sync.Map
)

// 当前的某个消息ID是否被屏蔽
func IsBlockedMessageByID(msgid int) bool {

	_, ok := blockedMsgByID.Load(msgid)

	return ok
}

// 按指定规则(或消息名)屏蔽消息日志
func BlockMessageLog(nameRule string) (err error, matchCount int) {

	err = cellnet.MessageMetaVisit(nameRule, func(meta *cellnet.MessageMeta) bool {

		blockedMsgByID.Store(int(meta.ID), meta)
		matchCount++

		return true
	})

	return
}
