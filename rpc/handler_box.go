package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type BoxHandler struct {
	cellnet.BaseEventHandler
}

func (self *BoxHandler) Call(ev *cellnet.SessionEvent) error {

	msgID := ev.MsgID

	// 来自encode之后的消息
	ev.FromMessage(&coredef.RemoteCallACK{
		MsgID: ev.MsgID,
		Data:  ev.Data,
	})

	// 消息ID不能用RemoteCallACK, 还是用消息本体
	ev.MsgID = msgID

	return self.CallNext(ev)
}

func NewBoxHandler() cellnet.EventHandler {

	return &BoxHandler{}

}
