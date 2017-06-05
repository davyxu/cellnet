package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler(useMsgLog bool) []cellnet.EventHandler {

	return cellnet.HandlerLink(cellnet.StaticEncodePacketHandler(),
		cellnet.HandlerOptional{StaticMsgLogHandler(), useMsgLog},
		StaticWritePacketHandler(),
	)

}

// socket.ReadPacketHandler -> socket.MsgLogHandler ->  socket.DecodePacketHandler -> cellnet.DispatcherHandler
func BuildRecvHandler(useMsgLog bool, recvHandler ...cellnet.EventHandler) []cellnet.EventHandler {

	return cellnet.HandlerLink(StaticReadPacketHandler(),
		cellnet.HandlerOptional{StaticMsgLogHandler(), useMsgLog},
		cellnet.StaticDecodePacketHandler(),
		recvHandler,
	)

}
