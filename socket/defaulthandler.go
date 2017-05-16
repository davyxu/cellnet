package socket

import "github.com/davyxu/cellnet"

// socket.EncodePacketHandler -> socket.MsgLogHandler -> socket.WritePacketHandler
func BuildSendHandler(useMsgLog bool) cellnet.EventHandler {

	if useMsgLog {

		return cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), NewMsgLogHandler(), NewWritePacketHandler())
	} else {

		return cellnet.LinkHandler(cellnet.NewEncodePacketHandler(), NewWritePacketHandler())
	}

}

// socket.ReadPacketHandler -> socket.MsgLogHandler -> socket.DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func BuildRecvHandler(useMsgLog bool, dispatcher *cellnet.DispatcherHandler, q cellnet.EventQueue) cellnet.EventHandler {

	if useMsgLog {

		return cellnet.LinkHandler(NewReadPacketHandler(), NewMsgLogHandler(), cellnet.NewDecodePacketHandler(), dispatcher)

	} else {
		return cellnet.LinkHandler(NewReadPacketHandler(), cellnet.NewDecodePacketHandler(), dispatcher)
	}

}
