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

type RPCHooker struct {
}

func (self RPCHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return
	}

	msg, meta, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return
	}

	peerInfo := inputEvent.Session().Peer().(interface {
		Name() string
	})

	log.Debugf("#rpc recv(%s)@%d %s(%d) | %s",
		peerInfo.Name(),
		inputEvent.Session().ID(),
		meta.TypeName(),
		meta.ID,
		cellnet.MessageToString(msg))

	poster := inputEvent.Session().Peer().(peer.MessagePoster)

	switch inputEvent.Message().(type) {
	case *RemoteCallREQ: // 服务端收到客户端的请求

		poster.PostEvent(&RecvMsgEvent{inputEvent.Session(), msg, rpcMsg.GetCallID()})

	case *RemoteCallACK: // 客户端收到服务器的回应
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(msg)
		}
	}

	return inputEvent
}

func (self RPCHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return
	}

	msg, meta, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return
	}

	peerInfo := inputEvent.Session().Peer().(interface {
		Name() string
	})

	log.Debugf("#rpc send(%s)@%d %s(%d) | %s",
		peerInfo.Name(),
		inputEvent.Session().ID(),
		meta.TypeName(),
		meta.ID,
		cellnet.MessageToString(msg))

	return inputEvent
}

func init() {
	msglog.BlockMessageLog("rpc.RemoteCallREQ")
	msglog.BlockMessageLog("rpc.RemoteCallACK")
}
