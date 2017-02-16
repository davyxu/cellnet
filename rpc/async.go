package rpc

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(ud interface{}, reqMsg interface{}, userCallback interface{}) error {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return err
	}

	recv, send := p.GetHandler()

	installSendHandler(p, send)

	if err := installAsyncRecvHandler(p, recv, reqMsg, userCallback); err != nil {
		return err
	}

	ses.Send(reqMsg)

	return nil
}

// 安装异步的接收回调
func installAsyncRecvHandler(p cellnet.Peer, recv cellnet.EventHandler, reqMsg interface{}, userCallback interface{}) error {

	funcType := reflect.TypeOf(userCallback)

	msgName := cellnet.MessageFullName(funcType.In(0))
	meta := cellnet.MessageMetaByName(msgName)
	if meta == nil {
		return ErrReplayMessageNotFound
	}

	// RPC消息不能用于普通消息
	// TODO 客户端请求时, 可以注册多个, 在处理完成时, 删除回调
	if p.CountByID(int(meta.ID)) == 0 {

		hl := cellnet.LinkHandler(
			socket.NewDecodePacketHandler(metaWrapper), // RemoteCall的Meta
			NewUnboxHandler(),
			socket.NewDecodePacketHandler(meta),
			NewReflectCallHandler(userCallback),
		)

		p.AddHandler(int(meta.ID), hl)
	}

	return nil
}
