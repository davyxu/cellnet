package rpc

import "github.com/davyxu/cellnet"

// Dispatcher( RemoteCallACK ) -> Unbox -> RPCMatch( RPCMsgID )-> Decode -> QueuePost -> Tail
func installACKHandler(p cellnet.Peer, id int, tailHandler cellnet.EventHandler) {

	var dispatcher *cellnet.DispatcherHandler

	raw := p.GetHandlerByIndex(int(metaWrapper.ID), 0)
	if raw == nil {
		// rpc服务端收到消息时, 用定制的handler返回消息, 而不是peer默认的
		raw = cellnet.LinkHandler(NewUnboxHandler(getSendHandler()), cellnet.NewDispatcherHandler())

		p.AddHandler(int(metaWrapper.ID), raw)
	}

	dispatcher = raw.Next().(*cellnet.DispatcherHandler)

	dispatcher.AddHandler(id, cellnet.LinkHandler(
		cellnet.NewDecodePacketHandler(),
		cellnet.NewQueuePostHandler(p.Queue()),
		tailHandler,
	))
}
