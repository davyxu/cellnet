package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

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

var metaWrapper = cellnet.MessageMetaByName("gamedef.RemoteCallACK")
