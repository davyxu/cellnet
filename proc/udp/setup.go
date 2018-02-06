package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
)

type MessageProc struct {
}

func (MessageProc) OnRecvMessage(ses cellnet.BaseSession) (msg interface{}, err error) {

	data := ses.Raw().(udp.DataReader).ReadData()

	return RecvLTVPacket(data)
}

func (MessageProc) OnSendMessage(ses cellnet.BaseSession, msg interface{}) error {

	writer := ses.(udp.DataWriter)

	return SendLTVPacket(writer, msg)
}

func init() {

	msgProc := new(MessageProc)
	msgLogger := new(msglog.LogHooker)

	proc.RegisterEventProcessor("udp.ltv", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHooker(msgLogger)
		initor.SetEventHandler(cellnet.UserMessageHandlerQueued(userHandler))

	})
}
