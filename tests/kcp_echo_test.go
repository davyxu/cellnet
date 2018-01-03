package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/kcpproc"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	"github.com/davyxu/cellnet/util"
	"testing"
	"time"
)

const kcpEcho_Address = "127.0.0.1:7903"

var kcpEcho_Signal *util.SignalTester

var kcpEcho_Acceptor cellnet.Peer

func kcpEcho_StartServer() {

	kcpEcho_Acceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Acceptor",
		EventProcessor: "udp.kcp",
		PeerAddress:    kcpEcho_Address,
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

func kcpEcho_StartClient() {

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Connector",
		EventProcessor: "udp.kcp",
		PeerAddress:    kcpEcho_Address,
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

					kcpEcho_Signal.Done(1)

				case *comm.SessionClosed:
					fmt.Println("client error: ")
				}
			}

			return nil
		},
	}).Start()

	kcpEcho_Signal.WaitAndExpect("not recv data", 1)
}

func TestKCPEcho(t *testing.T) {

	kcpEcho_Signal = util.NewSignalTester(t)
	kcpEcho_Signal.SetTimeout(time.Hour)

	kcpEcho_StartServer()

	time.Sleep(time.Millisecond)

	kcpEcho_StartClient()

	kcpEcho_Acceptor.Stop()
}
