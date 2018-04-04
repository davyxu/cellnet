package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer"
)

type RemoteCallMsg interface {
	GetMsgID() uint16
	GetMsgData() []byte
	GetCallID() int64
}

func ProcInboundEvent(inputEvent cellnet.Event) (handled bool) {

	if _, ok := inputEvent.(*RecvMsgEvent); ok {
		return true
	}

	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false
	}

	msg, meta, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false
	}

	if log.IsDebugEnabled() {

		peerInfo := inputEvent.Session().Peer().(interface {
			Name() string
		})

		log.Debugf("#rpc.recv(%s)@%d %s(%d) | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			meta.TypeName(),
			meta.ID,
			cellnet.MessageToString(msg))
	}

	switch inputEvent.Message().(type) {
	case *RemoteCallREQ: // 服务端收到客户端的请求

		poster := inputEvent.Session().Peer().(peer.MessagePoster)
		poster.PostEvent(&RecvMsgEvent{inputEvent.Session(), msg, rpcMsg.GetCallID()})

		// 避免后续环节处理
		return true

	case *RemoteCallACK: // 客户端收到服务器的回应
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(msg)
		}

		return true
	}

	return false
}

func ProcOutboundEvent(inputEvent cellnet.Event) (handled bool) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false
	}

	msg, meta, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false
	}

	if log.IsDebugEnabled() {

		peerInfo := inputEvent.Session().Peer().(interface {
			Name() string
		})

		log.Debugf("#rpc.send(%s)@%d %s(%d) | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			meta.TypeName(),
			meta.ID,
			cellnet.MessageToString(msg))
	}

	// 避免后续环节处理

	return true
}

func init() {
	msglog.BlockMessageLog("rpc.RemoteCallREQ")
	msglog.BlockMessageLog("rpc.RemoteCallACK")
}
