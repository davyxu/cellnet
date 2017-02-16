package rpc

import (
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

func installRecvHandler(p cellnet.Peer, recv cellnet.EventHandler, args interface{}, userCallback func(*cellnet.SessionEvent)) {

	// 接收
	msgName := cellnet.MessageFullName(reflect.TypeOf(args))
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		panic("can not found rpc message:" + msgName)
	}

	// RPC消息只能被注册1个
	// TODO 客户端请求时, 可以注册多个, 在处理完成时, 删除回调
	if p.CountByID(int(meta.ID)) == 0 {

		p.AddHandler(int(meta.ID), buildRecvHandler(meta, userCallback, nil))
	}

}
