package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler(useMsgLog bool) []cellnet.EventHandler {

	return cellnet.HandlerLink(cellnet.EncodePacketHandler(),
		cellnet.HandlerOptional{MsgLogHandler(), useMsgLog},
		WritePacketHandler(),
	)

}

// socket.ReadPacketHandler -> socket.MsgLogHandler -> socket.DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func BuildRecvHandler(useMsgLog bool, dispatcher *cellnet.DispatcherHandler) []cellnet.EventHandler {

	return cellnet.HandlerLink(ReadPacketHandler(),
		cellnet.HandlerOptional{MsgLogHandler(), useMsgLog},
		cellnet.DecodePacketHandler(),
		dispatcher,
	)

}
