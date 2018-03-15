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

type KCPMessageTransmitter struct {
}

func (KCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

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

func (KCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	return mustKCPContext(ses).sendMessage(msg)
}

type kcpEventHooker struct {
}

func (self kcpEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch inputEvent.Message().(type) {
	case *cellnet.SessionInit:
		inputEvent.Session().(cellnet.ContextSet).SetContext(kcpTag, newContext(inputEvent.Session()))
	case *cellnet.SessionClosed:
		mustKCPContext(inputEvent.Session()).Close()
	}

	return inputEvent
}

func (self kcpEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	transmitter := new(KCPMessageTransmitter)
	hooker := new(kcpEventHooker)

	proc.RegisterEventProcessor("udp.kcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetEventTransmitter(transmitter)
		bundle.SetEventHooker(hooker)
		bundle.SetEventCallback(cellnet.NewQueuedEventCallback(userCallback))

	})
}
