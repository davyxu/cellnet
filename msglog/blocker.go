package msglog

import (
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	blockedMsgByID sync.Map
)

func IsBlockedMessageByID(msgid int) bool {

	_, ok := blockedMsgByID.Load(msgid)

	return ok
}

func BlockMessageLog(nameRule string) (err error, matchCount int) {

	err = cellnet.MessageMetaVisit(nameRule, func(meta *cellnet.MessageMeta) bool {

		blockedMsgByID.Store(int(meta.ID), meta)
		matchCount++

		return true
	})

	return
}
