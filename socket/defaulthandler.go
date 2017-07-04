package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler() []cellnet.EventHandler {

	return cellnet.HandlerLink(cellnet.StaticEncodePacketHandler(),
		cellnet.StaticMsgLogHandler(),
		StaticWritePacketHandler(),
	)

}

// socket.ReadPacketHandler -> socket.MsgLogHandler ->  socket.DecodePacketHandler -> cellnet.DispatcherHandler
func BuildRecvHandler(recvHandler ...cellnet.EventHandler) []cellnet.EventHandler {

	return cellnet.HandlerLink(StaticReadPacketHandler(),
		cellnet.StaticMsgLogHandler(),
		cellnet.StaticDecodePacketHandler(),
		recvHandler,
	)

}
