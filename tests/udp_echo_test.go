package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	_ "github.com/davyxu/cellnet/comm/udpproc"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpEcho_Address = "127.0.0.1:7901"

var udpEcho_Signal *util.SignalTester

var udpEcho_Acceptor cellnet.Peer

func udpEcho_StartServer() {

	udpEcho_Acceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Acceptor",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpEcho_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&TestEchoACK{
						Msg:   msg.Msg,
						Value: msg.Value,
					})
				}
			}

			return nil
		},
	}).Start()

}

func udpEcho_StartClient() {

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Connector",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpEcho_Address,
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

					udpEcho_Signal.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	udpEcho_Signal.WaitAndExpect("not recv data", 1)
}

func TestUDPEcho(t *testing.T) {

	udpEcho_Signal = util.NewSignalTester(t)

	udpEcho_StartServer()

	udpEcho_StartClient()

	udpEcho_Acceptor.Stop()
}
