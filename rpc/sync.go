package rpc

import "github.com/davyxu/cellnet"

// 发出请求, 接收到服务器返回后才返回, ud: peer/session,   reqMsg:请求用的消息, ackMsgName: 返回消息类型名, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, ackMsgName string) (interface{}, error) {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	recv, send := p.GetHandler()

	installSendHandler(p, send)

	ret := make(chan interface{})

	if err := installSyncRecvHandler(p, recv, reqMsg, ackMsgName, ret); err != nil {
		return nil, err
	}

	ses.Send(reqMsg)

	return <-ret, nil
}

// 安装同步的接收回调
func installSyncRecvHandler(p cellnet.Peer, recv cellnet.EventHandler, reqMsg interface{}, msgName string, retChan chan interface{}) error {

	meta := cellnet.MessageMetaByName(msgName)
	if meta == nil {
		return ErrReplayMessageNotFound
	}

	// RPC消息只能被注册1个
	// TODO 客户端请求时, 可以注册多个, 在处理完成时, 删除回调
	if p.CountByID(int(meta.ID)) == 0 {

		hl := cellnet.LinkHandler(
			cellnet.NewDecodePacketHandler(metaWrapper), // RemoteCall的Meta
			NewUnboxHandler(nil),
			cellnet.NewDecodePacketHandler(meta),
			NewRetChanHandler(retChan),
		)

		p.AddHandler(int(meta.ID), hl)
	}

	return nil
}
