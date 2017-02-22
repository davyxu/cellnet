package rpc

import "github.com/davyxu/cellnet"

// Dispatcher( RemoteCallACK ) -> Unbox -> RPCMatch( RPCMsgID )-> Decode -> QueuePost -> Tail
func installReqHandler(p cellnet.Peer, id int, tailHandler cellnet.EventHandler) (rpcID int32) {

	var rpcDispatcher *RPCMatchHandler

	raw := p.GetHandlerByIndex(int(metaWrapper.ID), 0)
	if raw == nil {
		// rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		raw = cellnet.LinkHandler(NewUnboxHandler(getSendHandler()), NewRPCMatchHandler())

		p.AddHandler(int(metaWrapper.ID), raw)
	}

	rpcDispatcher = raw.Next().(*RPCMatchHandler)

	rpcID = int32(rpcDispatcher.AddHandler(id, cellnet.LinkHandler(
		cellnet.NewDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue()),
		tailHandler,
	)))

	tailHandler.SetNext(cellnet.NewCallbackHandler(func(ev *cellnet.SessionEvent) {

		rpcDispatcher.RemoveHandler(id, int(rpcID))
	}))
	// 移除

	return rpcID

}
