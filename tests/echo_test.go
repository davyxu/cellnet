package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/packet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const echoAddress = "127.0.0.1:7701"

var echoSignal *util.SignalTester

var echoAcceptor cellnet.Peer

func OnEchoServerEvent(raw cellnet.EventParam) cellnet.EventResult {
	switch ev := raw.(type) {
	case socket.AcceptedEvent:
		fmt.Println("server accepted")
	case packet.RecvMsgEvent:

		msg := ev.Msg.(*proto.TestEchoACK)

		fmt.Printf("server recv %+v\n", msg)

		ev.Ses.Send(&proto.TestEchoACK{
			Msg:   msg.Msg,
			Value: msg.Value,
		})
	case socket.SessionClosedEvent:
		fmt.Println("server error: ", ev.Error)
	}

	return nil
}

func EchoServer() {
	queue := cellnet.NewEventQueue()

	echoAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tcp.Acceptor",
		Queue:       queue,
		PeerAddress: echoAddress,
		PeerName:    "server",
		Event:       packet.ProcTLVPacket(msglog.ProcMsgLog(OnEchoServerEvent)),
	}).Start()

	queue.StartLoop()
}

func OnEchoClientEvent(raw cellnet.EventParam) cellnet.EventResult {
	switch ev := raw.(type) {
	case socket.ConnectedEvent:
		fmt.Println("client connected")
		ev.Ses.Send(&proto.TestEchoACK{
			Msg:   "hello",
			Value: 1234,
		})
	case packet.RecvMsgEvent:

		msg := ev.Msg.(*proto.TestEchoACK)

		fmt.Printf("client recv %+v\n", msg)

		echoSignal.Done(1)

	case socket.SessionClosedEvent:
		fmt.Println("client error: ", ev.Error)
	}

	return nil
}

func EchoClient() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tcp.Connector",
		Queue:       queue,
		PeerAddress: echoAddress,
		PeerName:    "client",
		Event:       packet.ProcTLVPacket(msglog.ProcMsgLog(OnEchoClientEvent)),
	}).Start()

	queue.StartLoop()

	echoSignal.WaitAndExpect("not recv data", 1)
}

func TestEcho(t *testing.T) {

	echoSignal = util.NewSignalTester(t)

	EchoServer()

	EchoClient()

	echoAcceptor.Stop()
}
