package kcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/msglog"
)

const kcpTag = "kcp"

func mustKCPContext(ses cellnet.Session) (ctx *kcpContext) {
	if ses.(cellnet.PropertySet).GetProperty(kcpTag, &ctx) {
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

	recvingData = false
	return
}

func (MessageProc) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	return mustKCPContext(ses).sendMessage(msg)
}

type udpEventHooker struct {
	logger msglog.LogHooker
}

func (self udpEventHooker) OnInboundEvent(ev cellnet.Event) {

	self.logger.OnInboundEvent(ev)

	switch ev.Message().(type) {
	case *cellnet.SessionInit:
		ev.Session().(cellnet.PropertySet).SetProperty(kcpTag, newContext(ev.Session()))
	case *cellnet.SessionClosed:
		mustKCPContext(ev.Session()).Close()
	}
}

func (self udpEventHooker) OnOutboundEvent(ev cellnet.Event) {

	self.logger.OnOutboundEvent(ev)
}

func init() {

	msgProc := new(MessageProc)
	msgHooker := new(udpEventHooker)

	proc.RegisterEventProcessor("udp.kcp.ltv", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHooker(msgHooker)
		initor.SetEventHandler(cellnet.UserMessageHandlerQueued(userHandler))

	})
}
