package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/tcppeer"
	_ "github.com/davyxu/cellnet/comm/tcpproc"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const tcpEchoAddress = "127.0.0.1:7701"

var tcpEchoSignal *util.SignalTester

var tcpEchoAcceptor cellnet.Peer

func StartTCPEchoServer() {
	queue := cellnet.NewEventQueue()

	tcpEchoAcceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    tcpEchoAddress,
		PeerName:       "server",
		InboundEvent: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionAccepted:
					fmt.Println("server accepted")
				case *TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&TestEchoACK{
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

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Connector",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    tcpEchoAddress,
		PeerName:       "client",
		InboundEvent: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *comm.SessionConnected:
					fmt.Println("client connected")
					ev.Ses.Send(&TestEchoACK{
						Msg:   "hello",
						Value: 1234,
					})
				case *TestEchoACK:

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
