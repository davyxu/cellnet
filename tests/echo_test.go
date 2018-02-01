package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/peer/udp"
	_ "github.com/davyxu/cellnet/proc/kcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
	_ "github.com/davyxu/cellnet/proc/udp"
	"github.com/davyxu/cellnet/util"
	"testing"
)

type echoContext struct {
	Address   string
	Protocol  string
	Processor string
	Tester    *util.SignalTester
	Acceptor  cellnet.Peer
}

var (
	echoContexts = []*echoContext{
		{
			Address:   "127.0.0.1:7701",
			Protocol:  "tcp",
			Processor: "tcp.ltv",
		},
		{
			Address:   "127.0.0.1:7702",
			Protocol:  "udp",
			Processor: "udp.ltv",
		},
		{
			Address:   "127.0.0.1:7703",
			Protocol:  "udp",
			Processor: "udp.kcp",
		},
	}
)

func echo_StartServer(context *echoContext) {
	queue := cellnet.NewEventQueue()

	context.Acceptor = peer.CreatePeer(peer.CommunicateConfig{
		PeerType:       context.Protocol + ".Acceptor",
		EventProcessor: context.Processor,
		UserQueue:      queue,
		PeerAddress:    context.Address,
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

func echo_StartClient(context *echoContext) {
	queue := cellnet.NewEventQueue()

	peer.CreatePeer(peer.CommunicateConfig{
		PeerType:       context.Protocol + ".Connector",
		EventProcessor: context.Processor,
		UserQueue:      queue,
		PeerAddress:    context.Address,
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

					context.Tester.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	context.Tester.WaitAndExpect("not recv data", 1)
}

func runEcho(t *testing.T, index int) {

	ctx := echoContexts[index]

	ctx.Tester = util.NewSignalTester(t)

	echo_StartServer(ctx)

	echo_StartClient(ctx)

	ctx.Acceptor.Stop()
}

func TestEchoTCP(t *testing.T) {

	runEcho(t, 0)
}

func TestEchoUDP(t *testing.T) {

	runEcho(t, 1)
}

func TestEchoKCP(t *testing.T) {

	runEcho(t, 2)
}
