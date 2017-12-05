package rpc

import (
	"github.com/davyxu/cellnet"
)

func sendRequest(ses cellnet.Session, msg interface{}, callid int64) {

	data, msgid, _ := cellnet.EncodeMessage(msg)

	evFunc := ses.Peer().EventFunc()
	if evFunc != nil {

		evFunc(cellnet.SendEvent{ses, &RemoteCallREQ{
			MsgID:  msgid,
			Data:   data,
			CallID: callid,
		}})
	}
}
