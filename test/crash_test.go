package tests

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellmsglog "github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/tcp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRecvCrash(t *testing.T) {

	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	tcp.TestEnableRecvPanic = true
	acc := tcp.NewAcceptor()
	acc.CapturePanic = true
	acc.Recv = tcp.RecvMessage
	acc.Send = tcp.SendMessage
	acc.Inbound = func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent) {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionClosed:
			signal.Done(msg.Err.Error())
		}
		return input
	}

	assert.NoError(t, acc.ListenAndAccept(tcpListen), "ListenAndAccept failed")

	conn := tcp.NewConnector()
	conn.CapturePanic = true
	conn.Recv = tcp.RecvMessage
	conn.Send = tcp.SendMessage
	conn.Inbound = func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent) {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionClosed:
			signal.Done(msg.Err.Error())
		}
		return input
	}
	conn.AsyncConnect(tcpListen)
	signal.WaitAll("protect failed", "recv panic: emulate recv crash", "recv panic: emulate recv crash")
}

func TestSendCrash(t *testing.T) {

	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	tcp.TestEnableSendPanic = true
	acc := tcp.NewAcceptor()
	acc.CapturePanic = true
	acc.Recv = tcp.RecvMessage
	acc.Send = tcp.SendMessage
	acc.Inbound = func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent) {
		cellmsglog.RecvLogger(input)
		return input
	}

	assert.NoError(t, acc.ListenAndAccept(tcpListen), "ListenAndAccept failed")

	conn := tcp.NewConnector()
	conn.CapturePanic = true
	conn.Recv = tcp.RecvMessage
	conn.Send = tcp.SendMessage
	conn.Inbound = func(input *cellevent.RecvMsgEvent) (output *cellevent.RecvMsgEvent) {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionConnected:
			input.Ses.Send(&TestEchoACK{
				Msg: "hello",
			})
		case *cellevent.SessionClosed:
			signal.Done(msg.Err.Error())
		}
		return input
	}

	tcp.OnSendCrash = func(raw interface{}) {
		signal.Done(raw.(string))
	}
	conn.AsyncConnect(tcpListen)
	signal.WaitAll("protect failed", "emulate send crash")
}
