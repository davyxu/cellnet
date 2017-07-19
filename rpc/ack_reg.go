package rpc

import (
	"github.com/davyxu/cellnet"
)

func ChainSend() *cellnet.HandlerChain {

	return cellnet.NewHandlerChain(
		cellnet.StaticEncodePacketHandler(),
		NewBoxHandler(),
	)
}

// 服务器端响应RPC消息
func RegisterMessage(p cellnet.Peer, msgName string, userCallback func(ev *cellnet.Event)) {

	if p == nil {
		return
	}

	msgMeta := cellnet.MessageMetaByName(msgName)

	p.AddChainRecv(cellnet.NewHandlerChain(
		cellnet.NewMatchMsgIDHandler(MetaCall.ID),
		cellnet.StaticDecodePacketHandler(),
		NewUnboxHandler(ChainSend()),
		cellnet.NewMatchMsgIDHandler(msgMeta.ID),
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue(), cellnet.NewCallbackHandler(userCallback)),
	))

}
