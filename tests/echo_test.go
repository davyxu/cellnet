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

const testAddress = "127.0.0.1:7701"

var echoSignal *util.SignalTester

var echoAcceptor cellnet.Peer

func onServerEvent(raw cellnet.EventParam) cellnet.EventResult {
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

func server() {
	queue := cellnet.NewEventQueue()

	echoAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerTypeName: "tcp.Acceptor",
		Queue:        queue,
		PeerAddress:  testAddress,
		PeerName:     "server",
		Event:        packet.ProcTLVPacket(msglog.ProcMsgLog(onServerEvent)),
	}).Start()

	queue.StartLoop()
}

func onClientEvent(raw cellnet.EventParam) cellnet.EventResult {
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

func client() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerTypeName: "tcp.Connector",
		Queue:        queue,
		PeerAddress:  testAddress,
		PeerName:     "client",
		Event:        packet.ProcTLVPacket(msglog.ProcMsgLog(onClientEvent)),
	}).Start()

	queue.StartLoop()

	echoSignal.WaitAndExpect("not recv data", 1)
}

func TestEcho(t *testing.T) {

	echoSignal = util.NewSignalTester(t)

	server()

	client()

	echoAcceptor.Stop()
}
