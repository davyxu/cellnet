package msglog

import (
	"errors"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	msgMetaByID sync.Map
)

func IsBlockedMessageByID(msgid int) bool {

	_, ok := msgMetaByID.Load(msgid)

	return ok
}

var (
	ErrMessageNotFound = errors.New("msg not exists")
)

func BlockMessageLog(msgName string) error {
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		return ErrMessageNotFound
	}

	msgMetaByID.Store(int(meta.ID), meta)

	return nil
}
