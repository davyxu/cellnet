package socket

import "github.com/davyxu/cellnet"

type WritePacketHandler struct {
	cellnet.BaseEventHandler
}

func (self *WritePacketHandler) Call(ev *cellnet.SessionEvent) error {

	rawSes := ev.Ses.(*SocketSession)
	rawSes.sendList.Add(ev)

	return nil
}

func NewWritePacketHandler() cellnet.EventHandler {

	return &WritePacketHandler{}

}
