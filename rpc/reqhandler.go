package rpc

import "github.com/davyxu/cellnet"

// 请求方接收消息
// Dispatcher( RemoteCallACK ) -> Unbox -> RPCMatch( RPCMsgID )-> Decode -> QueuePost -> Tail
func buildRecvHandler(p cellnet.Peer, msgName string, tailHandler cellnet.EventHandler) (rpcID int32, err error) {

	meta := cellnet.MessageMetaByName(msgName)
	if meta == nil {
		return -1, ErrReplayMessageNotFound
	}

	id := int(meta.ID)

	var rpcDispatcher *RPCMatchHandler

	raw := p.GetHandlerByIndex(int(MetaCall.ID), 0)
	if raw == nil {
		// rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		raw = cellnet.HandlerLink(NewUnboxHandler(getSendHandler()), NewRPCMatchHandler())

		p.AddHandler(int(MetaCall.ID), raw)
	}

	rpcDispatcher = raw[len(raw)-1].(*RPCMatchHandler)

	rpcID = int32(rpcDispatcher.AddHandler(id, cellnet.HandlerLink(
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue(), tailHandler, cellnet.NewCallbackHandler(func(ev *cellnet.SessionEvent) {

			rpcDispatcher.RemoveHandler(id, int(rpcID))
		})),
	)))

	return rpcID, nil

}
