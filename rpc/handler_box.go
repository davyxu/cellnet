package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/coredef"
)

type BoxHandler struct {
	cellnet.BaseEventHandler
}

func (self *BoxHandler) Call(ev *cellnet.SessionEvent) {

	// 来自encode之后的消息
	ev.FromMessage(&coredef.RemoteCallACK{
		MsgID: ev.MsgID,
		Data:  ev.Data,
	})

}

func NewBoxHandler() cellnet.EventHandler {

	return &BoxHandler{}

}
