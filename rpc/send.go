package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

func sendRequest(ses cellnet.Session, msg interface{}, callid int64) {

	data, msgid, _ := cellnet.EncodeMessage(msg)

	evFunc := ses.Peer().EventFunc()
	if evFunc != nil {

		evFunc(socket.SendEvent{ses, &RemoteCallREQ{
			MsgID:  msgid,
			Data:   data,
			CallID: callid,
		}})
	}
}
