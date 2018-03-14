package msglog

import (
	"errors"
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

var (
	ErrMessageNotFound = errors.New("msg not exists")
)

func BlockMessageLog(msgName string) error {
	meta := cellnet.MessageMetaByFullName(msgName)

	if meta == nil {
		return ErrMessageNotFound
	}

	blockedMsgByID.Store(int(meta.ID), meta)

	return nil
}
