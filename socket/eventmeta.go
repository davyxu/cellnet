package socket

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/proto/gamedef"
)

var (
	Meta_SessionConnected     = cellnet.MessageMetaByName("gamedef.SessionConnected")
	Meta_SessionClosed        = cellnet.MessageMetaByName("gamedef.SessionClosed")
	Meta_SessionAccepted      = cellnet.MessageMetaByName("gamedef.SessionAccepted")
	Meta_SessionAcceptFailed  = cellnet.MessageMetaByName("gamedef.SessionAcceptFailed")
	Meta_SessionConnectFailed = cellnet.MessageMetaByName("gamedef.SessionConnectFailed")
)
