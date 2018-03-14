package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/rpc"
	"io"
)

type MessageProc struct {
}

func (MessageProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	msg, err = RecvLTVPacket(reader)

	msglog.WriteRecvLogger(log, "udp", ses, msg)

	return
}

func (MessageProc) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	return SendLTVPacket(writer, msg)
}

type rpcEventHooker struct {
	rpc.RPCHooker
}

func (self rpcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	self.RPCHooker.OnInboundEvent(inputEvent)

	return inputEvent
}

func (self rpcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	self.RPCHooker.OnOutboundEvent(inputEvent)

	return inputEvent
}

func init() {

	msgProc := new(MessageProc)
	msgLogger := new(rpcEventHooker)

	proc.RegisterEventProcessor("tcp.ltv", func(initor proc.ProcessorBundleInitor, userHandler cellnet.UserMessageHandler) {

		initor.SetEventProcessor(msgProc)
		initor.SetEventHooker(msgLogger)
		initor.SetEventHandler(cellnet.UserMessageHandlerQueued(userHandler))

	})
}
