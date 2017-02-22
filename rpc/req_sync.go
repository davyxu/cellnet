package rpc

import "github.com/davyxu/cellnet"

// 发出请求, 接收到服务器返回后才返回, ud: peer/session,   reqMsg:请求用的消息, ackMsgName: 返回消息类型名, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, ackMsgName string) (interface{}, error) {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	//_, send := p.GetHandler()

	//installSendHandler(p, send)

	ret := make(chan interface{})

	rpcid, err := installSyncRecvHandler(p, ackMsgName, ret)
	if err != nil {
		return nil, err
	}

	ev := cellnet.NewSessionEvent(cellnet.SessionEvent_Send, ses)
	ev.TransmitTag = rpcid
	ev.Msg = reqMsg
	ses.RawSend(getSendHandler(), ev)

	return <-ret, nil
}

// 安装同步的接收回调
func installSyncRecvHandler(p cellnet.Peer, msgName string, retChan chan interface{}) (rpcID int32, err error) {

	meta := cellnet.MessageMetaByName(msgName)
	if meta == nil {
		return -1, ErrReplayMessageNotFound
	}

	rpcID = installReqHandler(p, int(meta.ID), NewRetChanHandler(retChan))

	return
}
