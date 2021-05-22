package tests

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellmsglog "github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/tcp"
	tcptransmit "github.com/davyxu/cellnet/transmit/tcp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	crashListen1 = "127.0.0.1:7702"
	crashListen2 = "127.0.0.1:7703"
)

func TestRecvCrash(t *testing.T) {

	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	tcptransmit.TestEnableRecvPanic = true

	defer func() {
		tcptransmit.TestEnableRecvPanic = false
	}()

	acc := tcp.NewAcceptor()
	acc.CapturePanic = true
	acc.Recv = tcptransmit.RecvMessage
	acc.Send = tcptransmit.SendMessage
	acc.Inbound = func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg) {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionClosed:
			signal.Done(msg.Err.Error())
		}
		return input
	}

	assert.NoError(t, acc.ListenAndAccept(crashListen1), "ListenAndAccept failed")

	conn := tcp.NewConnector()
	conn.CapturePanic = true
	conn.Recv = tcptransmit.RecvMessage
	conn.Send = tcptransmit.SendMessage
	conn.Inbound = func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg) {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionClosed:
			signal.Done(msg.Err.Error())
		}
		return input
	}
	conn.AsyncConnect(crashListen1)
	signal.WaitAll("protect failed", "recv panic: emulate recv crash", "recv panic: emulate recv crash")

}

func TestSendCrash(t *testing.T) {

	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	tcptransmit.TestEnableSendPanic = true

	defer func() {
		tcptransmit.TestEnableSendPanic = false
	}()

	acc := tcp.NewAcceptor()
	acc.CapturePanic = true
	acc.Recv = tcptransmit.RecvMessage
	acc.Send = tcptransmit.SendMessage
	acc.Inbound = func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg) {
		cellmsglog.RecvLogger(input)
		return input
	}

	assert.NoError(t, acc.ListenAndAccept(crashListen2), "ListenAndAccept failed")

	conn := tcp.NewConnector()
	conn.CapturePanic = true
	conn.Recv = tcptransmit.RecvMessage
	conn.Send = tcptransmit.SendMessage
	conn.Inbound = func(input *cellevent.RecvMsg) (output *cellevent.RecvMsg) {
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
	conn.AsyncConnect(crashListen2)
	signal.WaitAll("protect failed", "emulate send crash")
}
