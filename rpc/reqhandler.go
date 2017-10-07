package rpc

import "github.com/davyxu/cellnet"

// 请求方(rpc服务器）接收消息处理过程
func buildRecvHandler(p cellnet.Peer, msgName string, tailHandler cellnet.EventHandler) (rpcID int64, err error) {

	msgMeta := cellnet.MessageMetaByName(msgName)
	if msgMeta == nil {
		return -1, ErrReplayMessageNotFound
	}

	rpcID = p.AddChainRecv(cellnet.NewHandlerChain(
		cellnet.NewMatchMsgIDHandler(MetaCall.ID),
		cellnet.StaticDecodePacketHandler(),
		NewUnboxHandler(ChainSend()), // 将rpc的发包过程保存在Handler中，每次解包时，能随着Event透传到Send中
		cellnet.NewMatchMsgIDHandler(msgMeta.ID),
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue(), tailHandler, cellnet.NewCallbackHandler(func(ev *cellnet.Event) {
			// 处理完毕时，移除这个处理链
			p.RemoveChainRecv(rpcID)

		})),
	))

	return rpcID, nil

}
