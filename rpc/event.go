package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type RecvMsgEvent struct {
	ses    cellnet.Session
	Msg    interface{}
	callid int64
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Queue() cellnet.EventQueue {
	return self.ses.Peer().(interface {
		Queue() cellnet.EventQueue
	}).Queue()
}

func (self *RecvMsgEvent) Reply(msg interface{}) {

	data, meta, err := codec.EncodeMessage(msg, nil)

	if err != nil {
		log.Errorf("rpc reply message encode error: %s", err)
		return
	}

	self.ses.Send(&RemoteCallACK{
		MsgID:  uint32(meta.ID),
		Data:   data,
		CallID: self.callid,
	})
}
