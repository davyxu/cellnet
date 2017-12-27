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

const tcpEcho_Address = "127.0.0.1:7701"

var tcpEcho_Signal *util.SignalTester

var tcpEcho_Acceptor cellnet.Peer

func tcpEcho_StartServer() {
	queue := cellnet.NewEventQueue()

	tcpEcho_Acceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    tcpEcho_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

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

func tcpEcho_StartClient() {
	queue := cellnet.NewEventQueue()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Connector",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    tcpEcho_Address,
		PeerName:       "client",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

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

					tcpEcho_Signal.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	tcpEcho_Signal.WaitAndExpect("not recv data", 1)
}

func TestTCPEcho(t *testing.T) {

	tcpEcho_Signal = util.NewSignalTester(t)

	tcpEcho_StartServer()

	tcpEcho_StartClient()

	tcpEcho_Acceptor.Stop()
}
