package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/msglog"
)

type RemoteCallMsg interface {
	GetMsgID() uint16
	GetMsgData() []byte
	GetCallID() int64
}

func ResolveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool, err error) {

	if _, ok := inputEvent.(*RecvMsgEvent); ok {
		return inputEvent, false, nil
	}

	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return inputEvent, false, nil
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return inputEvent, false, err
	}

	if msglog.IsMsgLogValid(int(rpcMsg.GetMsgID())) {
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
		}, true, nil

	case *RemoteCallACK: // 客户端收到服务器的回应
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(userMsg)
		}

		return inputEvent, true, nil
	}

	return inputEvent, false, nil
}

func ResolveOutboundEvent(inputEvent cellnet.Event) (handled bool, err error) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false, nil
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false, err
	}

	if msglog.IsMsgLogValid(int(rpcMsg.GetMsgID())) {
		peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

		log.Debugf("#rpc.send(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			cellnet.MessageToString(userMsg))
	}

	// 避免后续环节处理

	return true, nil
}
