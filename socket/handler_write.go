package socket

import "github.com/davyxu/cellnet"

type WritePacketHandler struct {
}

func (self *WritePacketHandler) Call(ev *cellnet.SessionEvent) {

	rawSes := ev.Ses.(*SocketSession)
	rawSes.sendList.Add(ev)

}

var defaultWritePacketHandler = new(WritePacketHandler)

func StaticWritePacketHandler() cellnet.EventHandler {
	return defaultWritePacketHandler
}
