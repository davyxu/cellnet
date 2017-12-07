package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/tcppeer"
	"github.com/davyxu/cellnet/tcppkt"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const echoAddress = "127.0.0.1:7701"

var echoSignal *util.SignalTester

var echoAcceptor cellnet.Peer

func onEchoServerEvent(raw cellnet.EventParam) cellnet.EventResult {

	ev, ok := raw.(cellnet.RecvMsgEvent)
	if ok {
		switch msg := ev.Msg.(type) {
		case *tcppkt.SessionAccepted:
			fmt.Println("server accepted")
		case *proto.TestEchoACK:

			fmt.Printf("server recv %+v\n", msg)

			ev.Ses.Send(&proto.TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})

		case *tcppkt.SessionClosed:
			fmt.Println("server error: ")
		}
	}

	return nil
}

func StartEchoServer() {
	queue := cellnet.NewEventQueue()

	echoAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Acceptor",
		Queue:       queue,
		PeerAddress: echoAddress,
		PeerName:    "server",
		Event:       onEchoServerEvent,
	}).Start()

	queue.StartLoop()
}

func onEchoClientEvent(raw cellnet.EventParam) cellnet.EventResult {

	ev, ok := raw.(cellnet.RecvMsgEvent)
	if ok {
		switch msg := ev.Msg.(type) {
		case *tcppkt.SessionConnected:
			fmt.Println("client connected")
			ev.Ses.Send(&proto.TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *proto.TestEchoACK:

			fmt.Printf("client recv %+v\n", msg)

			echoSignal.Done(1)

		case *tcppkt.SessionClosed:
			fmt.Println("client error: ")
		}
	}

	return nil
}

func StartEchoClient() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Connector",
		Queue:       queue,
		PeerAddress: echoAddress,
		PeerName:    "client",
		Event:       onEchoClientEvent,
	}).Start()

	queue.StartLoop()

	echoSignal.WaitAndExpect("not recv data", 1)
}

func TestEcho(t *testing.T) {

	echoSignal = util.NewSignalTester(t)

	StartEchoServer()

	StartEchoClient()

	echoAcceptor.Stop()
}
