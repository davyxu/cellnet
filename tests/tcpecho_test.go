package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/tcppeer"
	_ "github.com/davyxu/cellnet/comm/tcppkt"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const tcpEchoAddress = "127.0.0.1:7701"

var tcpEchoSignal *util.SignalTester

var tcpEchoAcceptor cellnet.Peer

func StartTCPEchoServer() {
	queue := cellnet.NewEventQueue()

	tcpEchoAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Acceptor",
		Queue:       queue,
		PeerAddress: tcpEchoAddress,
		PeerName:    "server",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionAccepted:
					fmt.Println("server accepted")
				case *proto.TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&proto.TestEchoACK{
						Msg:   msg.Msg,
						Value: msg.Value,
					})

				case *comm.SessionClosed:
					fmt.Println("session closed: ", ev.Ses.ID())
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()
}

func StartTCPEchoClient() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Connector",
		Queue:       queue,
		PeerAddress: tcpEchoAddress,
		PeerName:    "client",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionConnected:
					fmt.Println("client connected")
					ev.Ses.Send(&proto.TestEchoACK{
						Msg:   "hello",
						Value: 1234,
					})
				case *proto.TestEchoACK:

					fmt.Printf("client recv %+v\n", msg)

					tcpEchoSignal.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	tcpEchoSignal.WaitAndExpect("not recv data", 1)
}

func TestTCPEcho(t *testing.T) {

	tcpEchoSignal = util.NewSignalTester(t)

	StartTCPEchoServer()

	StartTCPEchoClient()

	tcpEchoAcceptor.Stop()
}
