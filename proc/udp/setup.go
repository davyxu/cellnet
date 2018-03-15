package udp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
)

type MessageProc struct {
}

func (MessageProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	data := ses.Raw().(udp.DataReader).ReadData()

	msg, err = RecvLTVPacket(data)

	msglog.WriteRecvLogger(log, "udp", ses, msg)

	return
}

func (MessageProc) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer := ses.(udp.DataWriter)

	msglog.WriteSendLogger(log, "udp", ses, msg)

	return SendLTVPacket(writer, msg)
}

func init() {

	msgProc := new(MessageProc)

	proc.RegisterEventProcessor("udp.ltv", func(initor proc.ProcessorBundleSetter, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHandler(cellnet.UserMessageHandlerQueued(userHandler))

	})
}
