package socket

import "github.com/davyxu/cellnet"

type writePacketHandler struct {
}

func (self *writePacketHandler) Call(ev *cellnet.SessionEvent) {

	rawSes := ev.Ses.(*SocketSession)
	rawSes.sendList.Add(ev)

}

var defaultWritePacketHandler = new(writePacketHandler)

func WritePacketHandler() cellnet.EventHandler {

	return defaultWritePacketHandler

}
