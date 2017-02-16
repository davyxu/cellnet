package rpc

import (
	"errors"
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

var metaWrapper = cellnet.MessageMetaByName("gamedef.RemoteCallACK")

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler

// 注册连接消息
// TODO 消息被标记为rpc消息后, 不可注册到普通消息里
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	p.AddHandler(int(meta.ID), buildRecvHandler(meta, userCallback, nil))

	_, send := p.GetHandler()
	installSendHandler(p, send)
}

func buildRecvHandler(meta *cellnet.MessageMeta, userCallback func(ev *cellnet.SessionEvent), signalHandler cellnet.EventHandler) cellnet.EventHandler {
	return cellnet.LinkHandler(socket.NewDecodePacketHandler(metaWrapper), NewUnboxHandler(), socket.NewDecodePacketHandler(meta), socket.NewCallbackHandler(userCallback), signalHandler)
}

var ErrReplayMessageNotFound = errors.New("Reply message name not found")

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
			socket.NewDecodePacketHandler(metaWrapper), // RemoteCall的Meta
			NewUnboxHandler(),
			socket.NewDecodePacketHandler(meta),
			NewRetChanHandler(retChan),
		)

		p.AddHandler(int(meta.ID), hl)
	}

	return nil
}
