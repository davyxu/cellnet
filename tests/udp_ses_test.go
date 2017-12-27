package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	_ "github.com/davyxu/cellnet/comm/udppeer"
	_ "github.com/davyxu/cellnet/comm/udpproc"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpSes_Address = "127.0.0.1:7902"

func TestUDPClientPositiveClose(t *testing.T) {

	signal := util.NewSignalTester(t)

	acceptor := cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Acceptor",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpSes_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionClosed:
					signal.Done(2)
				}
			}

			return nil
		},
	}).Start()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Connector",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpSes_Address,
		PeerName:       "client",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionConnected:
					ev.Ses.Close()
				case *comm.SessionClosed:
					signal.Done(1)
				}
			}

			return nil
		},
	}).Start()

	signal.WaitAndExpect("client not closed", 1)
	signal.WaitAndExpect("server not recv closed", 2)

	acceptor.Stop()
}

func TestUDPServerPositiveClose(t *testing.T) {

	signal := util.NewSignalTester(t)

	acceptor := cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Acceptor",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpSes_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionAccepted:
					ev.Ses.Close()
					signal.Done(1)
				}
			}

			return nil
		},
	}).Start()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "udp.Connector",
		EventProcessor: "udp.ltv",
		PeerAddress:    udpSes_Address,
		PeerName:       "client",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionConnected:
					ev.Ses.Send(&TestEchoACK{
						Msg:   "hello",
						Value: 1234,
					})
				case *comm.SessionClosed:
					signal.Done(2)
				}
			}

			return nil
		},
	}).Start()

	signal.WaitAndExpect("server not accept", 1)
	signal.WaitAndExpect("client not recv closed closed", 2)

	acceptor.Stop()
}
