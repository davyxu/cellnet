package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler(useMsgLog bool) []cellnet.EventHandler {

	return cellnet.HandlerLink(cellnet.StaticEncodePacketHandler(),
		cellnet.HandlerOptional{useMsgLog, cellnet.StaticMsgLogHandler()},
		StaticWritePacketHandler(),
	)

}

// socket.ReadPacketHandler -> socket.MsgLogHandler ->  socket.DecodePacketHandler -> cellnet.DispatcherHandler
func BuildRecvHandler(useMsgLog bool, recvHandler ...cellnet.EventHandler) []cellnet.EventHandler {

	return cellnet.HandlerLink(StaticReadPacketHandler(),
		cellnet.HandlerOptional{useMsgLog, cellnet.StaticMsgLogHandler()},
		cellnet.StaticDecodePacketHandler(),
		recvHandler,
	)

}
