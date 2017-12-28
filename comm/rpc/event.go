package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
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

	self.ses.Send(&comm.RemoteCallACK{
		MsgID:  meta.ID,
		Data:   data,
		CallID: self.callid,
	})

}
