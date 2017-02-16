package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler

// rpc服务器端注册连接消息
// TODO 消息被标记为rpc消息后, 不可注册到普通消息里
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	p.AddHandler(int(meta.ID), cellnet.LinkHandler(
		socket.NewDecodePacketHandler(metaWrapper),
		NewUnboxHandler(),
		socket.NewDecodePacketHandler(meta),
		socket.NewCallbackHandler(userCallback),
	))

	_, send := p.GetHandler()
	installSendHandler(p, send)
}
