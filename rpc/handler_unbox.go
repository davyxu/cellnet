package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/binary/coredef"
)

type UnboxHandler struct {
	feedbackChain *cellnet.HandlerChain
}

func (self *UnboxHandler) Call(ev *cellnet.Event) {

	wrapper := ev.Msg.(*coredef.RemoteCallACK)

	ev.MsgID = wrapper.MsgID
	ev.Data = wrapper.Data

	// 服务器接收后, 发送时, 需要使用CallID
	ev.TransmitTag = wrapper.CallID

	// 方便在rpc消息接收后，使用ev.Send直接将反馈消息发给rpc请求方。发送过程是一个需要定制的过程
	ev.ChainSend = self.feedbackChain

}

func NewUnboxHandler(chain *cellnet.HandlerChain) cellnet.EventHandler {
	return &UnboxHandler{
		feedbackChain: chain,
	}

}
