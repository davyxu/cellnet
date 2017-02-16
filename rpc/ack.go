package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler

func buildSendHandler() cellnet.EventHandler {

	if socket.EnableMessageLog {

		return cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), socket.NewMsgLogHandler(), NewBoxHandler(), socket.NewWritePacketHandler())
	} else {

		return cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), NewBoxHandler(), socket.NewWritePacketHandler())
	}
}

// rpc服务器端注册连接消息
// TODO 消息被标记为rpc消息后, 不可注册到普通消息里
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	sendHandler := buildSendHandler()

	p.AddHandler(int(meta.ID), cellnet.LinkHandler(
		cellnet.NewDecodePacketHandler(metaWrapper),
		NewUnboxHandler(sendHandler), // rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		cellnet.NewDecodePacketHandler(meta),
		cellnet.NewCallbackHandler(userCallback),
	))
}
