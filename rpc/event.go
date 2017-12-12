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

	data, meta, err := cellnet.EncodeMessage(msg)

	if err != nil {
		log.Errorf("rpc reply message encode error: %s", err)
		return
	}

	self.ses.Send(&RemoteCallACK{
		MsgID:  meta.ID,
		Data:   data,
		CallID: self.callid,
	})

}
