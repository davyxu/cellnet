package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type RemoteCallMsg interface {
	GetMsgID() uint16
	GetMsgData() []byte
	GetCallID() int64
}

func ResolveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool) {

	if _, ok := inputEvent.(*RecvMsgEvent); ok {
		return inputEvent, false
	}

	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return inputEvent, false
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return inputEvent, false
	}

	if log.IsDebugEnabled() {

		peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

		log.Debugf("#rpc.recv(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			cellnet.MessageToString(userMsg))
	}

	switch inputEvent.Message().(type) {
	case *RemoteCallREQ: // 服务端收到客户端的请求

		return &RecvMsgEvent{
			inputEvent.Session(),
			userMsg,
			rpcMsg.GetCallID(),
		}, true

	case *RemoteCallACK: // 客户端收到服务器的回应
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(userMsg)
		}

		return inputEvent, true
	}

	return inputEvent, false
}

func ResolveOutboundEvent(inputEvent cellnet.Event) (handled bool) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false
	}

	if log.IsDebugEnabled() {

		peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

		log.Debugf("#rpc.send(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			cellnet.MessageToString(userMsg))
	}

	// 避免后续环节处理

	return true
}
