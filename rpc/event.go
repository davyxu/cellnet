package rpc

import (
	"github.com/davyxu/cellnet"
)

type RecvMsgEvent struct {
	ses    cellnet.Session
	Msg    interface{}
	callid int64
}

func (self *RecvMsgEvent) Queue() cellnet.EventQueue {
	return self.ses.Peer().EventQueue()
}

func (self *RecvMsgEvent) Reply(msg interface{}) {

	data, msgid, _ := cellnet.EncodeMessage(msg)

	self.ses.Send(&RemoteCallACK{
		MsgID:  msgid,
		Data:   data,
		CallID: self.callid,
	})
}
