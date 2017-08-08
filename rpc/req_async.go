package rpc

import (
	"github.com/davyxu/cellnet"
	"time"
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(sesOrPeer interface{}, reqMsg interface{}, ackMsgName string, timeout time.Duration, userCallback func(ev *cellnet.Event)) error {

	ses, p, err := getPeerSession(sesOrPeer)

	if err != nil {
		return err
	}

	rpcid, err := buildRecvHandler(p, ackMsgName, cellnet.NewCallbackHandler(userCallback))

	if err != nil {
		return err
	}

	time.AfterFunc(timeout, func() {
		p.RemoveChainRecv(rpcid)
	})

	// 发送RPC请求
	ev := cellnet.NewEvent(cellnet.Event_Send, ses)
	ev.TransmitTag = rpcid
	ev.Msg = reqMsg
	ev.ChainSend = ChainSend()
	ses.RawSend(ev)

	return nil
}
