package tcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/rpc"
	"github.com/davyxu/cellnet/util"
	"io"
)

type TCPMessageTransmitter struct {
}

func (TCPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	msg, err = util.RecvLTVPacket(reader)

	msglog.WriteRecvLogger(log, "tcp", ses, msg)

	return
}

func (TCPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	return util.SendLTVPacket(writer, msg)
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

	transmitter := new(TCPMessageTransmitter)
	hooker := new(rpcEventHooker)

	proc.RegisterEventProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetEventTransmitter(transmitter)
		bundle.SetEventHooker(hooker)
		bundle.SetEventCallback(cellnet.NewQueuedEventCallback(userCallback))

	})
}
