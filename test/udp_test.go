package tests

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellmsglog "github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	udpListen = "127.0.0.1:7701"
)

func startUDPServer(t *testing.T) {
	acc := udp.NewAcceptor()
	acc.Recv = udp.RecvMessage
	acc.Send = udp.SendMessage
	acc.Outbound = cellmsglog.SendLogger
	acc.Inbound = func(input *cellevent.RecvMsgEvent) *cellevent.RecvMsgEvent {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *TestEchoACK:
			input.Ses.Send(msg)
		}

		return input
	}

	assert.NoError(t, acc.ListenAndAccept(udpListen), "ListenAndAccept failed")
}

func TestUDPEcho(t *testing.T) {
	startUDPServer(t)

	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	conn := udp.NewConnector()
	conn.Recv = udp.RecvMessage
	conn.Send = udp.SendMessage
	conn.Outbound = cellmsglog.SendLogger
	conn.Inbound = func(input *cellevent.RecvMsgEvent) *cellevent.RecvMsgEvent {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *cellevent.SessionConnected:
			input.Ses.Send(&TestEchoACK{
				Msg: "hello",
			})
		case *TestEchoACK:
			t.Log(msg)
			signal.Done(msg.Msg)
		}

		return input
	}

	conn.AsyncConnect(udpListen)

	signal.WaitAll("echo not respond", "hello")
}
