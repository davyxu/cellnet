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
// Read-> Decode-> Dispatcher -> Unbox -> Dispatcher(RPC) -> Decode-> QueuePost-> Callback
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	sendHandler := buildSendHandler()

	raw := p.GetHandlerByID(int(metaWrapper.ID))
	if raw == nil {
		// rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		raw = NewUnboxHandler(sendHandler)
		p.AddHandler(int(metaWrapper.ID), raw)
	}

	var rpcDispatcher *cellnet.DispatcherHandler
	rawDispatcher := raw.Next()
	if rawDispatcher == nil {
		rpcDispatcher = cellnet.NewDispatcherHandler()
		raw.SetNext(rpcDispatcher)
	} else {
		rpcDispatcher = rawDispatcher.(*cellnet.DispatcherHandler)
	}

	rpcDispatcher.AddHandler(int(meta.ID), cellnet.LinkHandler(
		cellnet.NewDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue()),
		cellnet.NewCallbackHandler(userCallback),
	))
}
