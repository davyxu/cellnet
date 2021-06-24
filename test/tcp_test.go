package tests

import (
	cellevent "github.com/davyxu/cellnet/event"
	cellmsglog "github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/tcp"
	tcptransport "github.com/davyxu/cellnet/transport/tcp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	tcpListen = "127.0.0.1:7701"
)

func startTCPServer(t *testing.T) {
	acc := tcp.NewAcceptor()
	acc.OnRecv = tcptransport.RecvMessage
	acc.OnSend = tcptransport.SendMessage
	acc.OnOutbound = cellmsglog.SendLogger
	acc.OnInbound = func(input *cellevent.RecvMsg) *cellevent.RecvMsg {
		cellmsglog.RecvLogger(input)
		switch msg := input.Message().(type) {
		case *TestEchoACK:
			input.Ses.Send(msg)
		}

		return input
	}

	assert.NoError(t, acc.ListenAndAccept(tcpListen), "ListenAndAccept failed")
}

func startTCPClient(t *testing.T) {
	signal := NewSignalTester(t)
	signal.SetTimeout(time.Second)

	conn := tcp.NewConnector()
	conn.OnRecv = tcptransport.RecvMessage
	conn.OnSend = tcptransport.SendMessage
	conn.OnOutbound = cellmsglog.SendLogger

	conn.OnInbound = func(input *cellevent.RecvMsg) *cellevent.RecvMsg {
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

	conn.AsyncConnect(tcpListen)

	signal.WaitAll("echo not respond", "hello")
}

func TestTCPEcho(t *testing.T) {
	startTCPServer(t)
	startTCPClient(t)
}
