package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler

var (
	sendHandlerWithLog cellnet.EventHandler
	sendHandler        cellnet.EventHandler
)

func getSendHandler() cellnet.EventHandler {

	if socket.EnableMessageLog {

		if sendHandlerWithLog == nil {
			sendHandlerWithLog = cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), socket.NewMsgLogHandler(), NewBoxHandler(), socket.NewWritePacketHandler())
		}

		return sendHandlerWithLog
	} else {

		if sendHandler == nil {
			sendHandler = cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), NewBoxHandler(), socket.NewWritePacketHandler())
		}

		return sendHandler
	}
}

// rpc服务器端注册连接消息
// Read-> Decode-> Dispatcher -> Unbox -> Dispatcher(RPC) -> Decode-> QueuePost-> Callback
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.SessionEvent)) {

	meta := cellnet.MessageMetaByName(msgName)

	installACKHandler(p, int(meta.ID), cellnet.NewCallbackHandler(userCallback) )
}
