package rpc

import (
	"github.com/davyxu/cellnet"
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(sesOrPeer interface{}, reqMsg interface{}, ackMessageName string, userCallback func(ev *cellnet.SessionEvent)) error {

	ses, p, err := getPeerSession(sesOrPeer)

	if err != nil {
		return err
	}

	rpcid, err := installAsyncRecvHandler(p, ackMessageName, userCallback)

	if err != nil {
		return err
	}

	ev := cellnet.NewSessionEvent(cellnet.SessionEvent_Send, ses)
	ev.TransmitTag = rpcid
	ev.Msg = reqMsg
	ses.RawSend(getSendHandler(), ev)

	return nil
}

// 安装异步的接收回调
func installAsyncRecvHandler(p cellnet.Peer, ackMessageName string, userCallback func(ev *cellnet.SessionEvent)) (rpcID int32, err error) {

	meta := cellnet.MessageMetaByName(ackMessageName)

	if meta == nil {
		return -1, ErrReplayMessageNotFound
	}

	rpcID = installReqHandler(p, int(meta.ID), cellnet.NewCallbackHandler(userCallback))

	return
}
