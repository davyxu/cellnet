package socket

import (
	"github.com/davyxu/cellnet"
)

var (
	Meta_SessionConnected     = cellnet.MessageMetaByName("coredef.SessionConnected")
	Meta_SessionClosed        = cellnet.MessageMetaByName("coredef.SessionClosed")
	Meta_SessionAccepted      = cellnet.MessageMetaByName("coredef.SessionAccepted")
	Meta_SessionAcceptFailed  = cellnet.MessageMetaByName("coredef.SessionAcceptFailed")
	Meta_SessionConnectFailed = cellnet.MessageMetaByName("coredef.SessionConnectFailed")
)
