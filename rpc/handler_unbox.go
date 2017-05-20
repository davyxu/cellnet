package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type UnboxHandler struct {
	feedbackHandler []cellnet.EventHandler
}

func (self *UnboxHandler) Call(ev *cellnet.SessionEvent) {

	wrapper := ev.Msg.(*coredef.RemoteCallACK)

	ev.MsgID = wrapper.MsgID
	ev.Data = wrapper.Data

	// 服务器接收后, 发送时, 需要使用CallID
	ev.TransmitTag = wrapper.CallID

	ev.SendHandler = self.feedbackHandler

}

func NewUnboxHandler(feedbackHandler []cellnet.EventHandler) cellnet.EventHandler {
	return &UnboxHandler{
		feedbackHandler: feedbackHandler,
	}

}
