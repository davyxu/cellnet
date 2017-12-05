package rpc

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	ses    cellnet.Session
	Msg    interface{}
	callid int64
}

func (self *RecvMsgEvent) Reply(msg interface{}) {

	data, msgid, _ := cellnet.EncodeMessage(msg)

	evFunc := self.ses.Peer().EventFunc()
	if evFunc != nil {

		evFunc(cellnet.SendEvent{self.ses, &RemoteCallACK{
			MsgID:  msgid,
			Data:   data,
			CallID: self.callid,
		}})
	}
}
