package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler

var (
	sendHandler []cellnet.EventHandler
)

func getSendHandler() []cellnet.EventHandler {

	if sendHandler == nil {
		sendHandler = cellnet.HandlerLink(cellnet.StaticEncodePacketHandler(),
			cellnet.HandlerOptional{socket.StaticMsgLogHandler(), socket.EnableMessageLog},
			NewBoxHandler(),
			socket.StaticWritePacketHandler(),
		)
	}

	return sendHandler

}

// 服务器端响应RPC消息
// Read-> Decode-> Dispatcher -> Unbox -> Dispatcher(RPC) -> Decode-> QueuePost-> Callback
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	id := int(meta.ID)

	var dispatcher *cellnet.DispatcherHandler

	raw := p.GetHandlerByIndex(int(metaWrapper.ID), 0)
	if raw == nil {
		// rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		raw = cellnet.HandlerLink(NewUnboxHandler(getSendHandler()), cellnet.NewDispatcherHandler())

		p.AddHandler(int(metaWrapper.ID), raw)
	}

	dispatcher = raw[len(raw)-1].(*cellnet.DispatcherHandler)

	dispatcher.AddHandler(id, cellnet.HandlerLink(
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue(), cellnet.HandlerLink(cellnet.NewCallbackHandler(userCallback))),
	))
}
