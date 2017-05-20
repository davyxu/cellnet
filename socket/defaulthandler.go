package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler(useMsgLog bool) []cellnet.EventHandler {

	return cellnet.HandlerLink(cellnet.StaticEncodePacketHandler(),
		cellnet.HandlerOptional{StaticMsgLogHandler(), useMsgLog},
		StaticWritePacketHandler(),
	)

}

// socket.ReadPacketHandler -> socket.MsgLogHandler -> socket.DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func BuildRecvHandler(useMsgLog bool, dispatcher *cellnet.DispatcherHandler) []cellnet.EventHandler {

	return cellnet.HandlerLink(StaticReadPacketHandler(),
		cellnet.HandlerOptional{StaticMsgLogHandler(), useMsgLog},
		cellnet.StaticDecodePacketHandler(),
		dispatcher,
	)

}
