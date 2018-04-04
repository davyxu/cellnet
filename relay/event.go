package relay

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	Ses cellnet.Session
	Msg interface{}

	SessionID []int64
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func init() {
	// 使用hubpeer的服务，看不到上下行的中间消息
	//msglog.BlockMessageLog("gamedef.HubPubUpstreamACK")
	//msglog.BlockMessageLog("gamedef.HubPubDownstreamACK")
}
