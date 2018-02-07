package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/msglog"
	"github.com/davyxu/cellnet/proc/rpc"
	"io"
)

type MessageProc struct {
}

func (MessageProc) OnRecvMessage(ses cellnet.BaseSession) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	return RecvLTVPacket(reader)
}

func (MessageProc) OnSendMessage(ses cellnet.BaseSession, msg interface{}) error {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	return SendLTVPacket(writer, msg)
}

type rpcEventHooker struct {
	rpc.RPCHooker
	msglog.LogHooker
}

func (self rpcEventHooker) OnInboundEvent(ev cellnet.Event) {

	self.LogHooker.OnInboundEvent(ev)
	self.RPCHooker.OnInboundEvent(ev)

}

func (self rpcEventHooker) OnOutboundEvent(ev cellnet.Event) {

	self.LogHooker.OnOutboundEvent(ev)
	self.RPCHooker.OnOutboundEvent(ev)
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
