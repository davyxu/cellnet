package kcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
)

const kcpTag = "kcp"

func mustKCPContext(ses cellnet.Session) (ctx *kcpContext) {
	if ses.(cellnet.ContextSet).GetContext(kcpTag, &ctx) {
		return
	} else {
		panic("invalid kcp context")
	}
}

type MessageProc struct {
}

func (MessageProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	ctx := mustKCPContext(ses)

	var recvingData = true

	go func() {
		for recvingData {

			data := ses.Raw().(udp.DataReader).ReadData()

			ctx.input(data)
		}
	}()

	msg, err = ctx.RecvLTVPacket()

	msglog.WriteRecvLogger(log, "kcp", ctx.ses, msg)

	recvingData = false
	return
}

func (MessageProc) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	return mustKCPContext(ses).sendMessage(msg)
}

type udpEventHooker struct {
}

func (self udpEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch inputEvent.Message().(type) {
	case *cellnet.SessionInit:
		inputEvent.Session().(cellnet.ContextSet).SetContext(kcpTag, newContext(inputEvent.Session()))
	case *cellnet.SessionClosed:
		mustKCPContext(inputEvent.Session()).Close()
	}

	return inputEvent
}

func (self udpEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	msgProc := new(MessageProc)
	msgHooker := new(udpEventHooker)

	proc.RegisterEventProcessor("udp.kcp.ltv", func(initor proc.ProcessorBundleSetter, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHooker(msgHooker)
		initor.SetEventHandler(cellnet.UserMessageHandlerQueued(userHandler))

	})
}
