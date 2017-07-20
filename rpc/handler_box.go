package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/binary/coredef"
)

type BoxHandler struct {
}

func (self *BoxHandler) Call(ev *cellnet.Event) {

	// 来自encode之后的消息
	ev.FromMessage(&coredef.RemoteCallACK{
		MsgID:  ev.MsgID,
		Data:   ev.Data,
		CallID: ev.TransmitTag.(int64),
	})

}

func NewBoxHandler() cellnet.EventHandler {

	return &BoxHandler{}

}
