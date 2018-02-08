package http

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

type totalProc struct {
	MessageProc
	StaticFileProc
}

func (self totalProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	msg, err = self.MessageProc.OnRecvMessage(ses)

	if err == errNotHandled {
		msg, err = self.StaticFileProc.OnRecvMessage(ses)
	}

	return
}

func (self totalProc) OnSendMessage(ses cellnet.Session, raw interface{}) error {

	return self.MessageProc.OnSendMessage(ses, raw)
}

func init() {

	totalProc := new(totalProc)
	msgProc := new(MessageProc)
	fileProc := new(StaticFileProc)
	msgLogger := new(LogHooker)

	proc.RegisterEventProcessor("http", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(totalProc)
		initor.SetEventHooker(msgLogger)
		initor.SetEventHandler(userHandler)
	})

	proc.RegisterEventProcessor("httpmsg", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHooker(msgLogger)
		initor.SetEventHandler(userHandler)
	})

	proc.RegisterEventProcessor("httpfile", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(fileProc)
		initor.SetEventHooker(msgLogger)
		initor.SetEventHandler(userHandler)
	})
}
