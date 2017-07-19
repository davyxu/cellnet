package rpc

import "github.com/davyxu/cellnet"

// 请求方接收消息
func buildRecvHandler(p cellnet.Peer, msgName string, tailHandler cellnet.EventHandler) (rpcID int64, err error) {

	msgMeta := cellnet.MessageMetaByName(msgName)
	if msgMeta == nil {
		return -1, ErrReplayMessageNotFound
	}

	rpcID = p.AddChainRecv(cellnet.NewHandlerChain(
		cellnet.NewMatchMsgIDHandler(MetaCall.ID),
		cellnet.StaticDecodePacketHandler(),
		NewUnboxHandler(ChainSend()),
		cellnet.NewMatchMsgIDHandler(msgMeta.ID),
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue(), tailHandler, cellnet.NewCallbackHandler(func(ev *cellnet.Event) {

			p.RemoveChainRecv(rpcID)

		})),
	))

	return rpcID, nil

}
