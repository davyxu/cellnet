package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
)

type UnboxHandler struct {
	cellnet.BaseEventHandler
}

func (self *UnboxHandler) Call(ev *cellnet.SessionEvent) error {

	wrapper := ev.Msg.(*gamedef.RemoteCallACK)

	ev.MsgID = wrapper.MsgID
	ev.Meta = cellnet.MessageMetaByID(wrapper.MsgID)
	ev.Data = wrapper.Data

	return self.CallNext(ev)
}

func NewUnboxHandler() cellnet.EventHandler {
	return &UnboxHandler{}

}
