package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/udp"
	"github.com/davyxu/cellnet/util"
	"testing"
)

const udpSes_Address = "127.0.0.1:10914"

func TestUDPClientPositiveClose(t *testing.T) {

	signal := util.NewSignalTester(t)

	acc := peer.NewGenericPeer("udp.Acceptor", "server", udpSes_Address, nil)

	proc.BindProcessorHandler(acc, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:
			signal.Done(2)
		}
	})

	acc.Start()

	connector := peer.NewGenericPeer("udp.Connector", "client", udpSes_Address, nil)

	proc.BindProcessorHandler(connector, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			ev.Session().Close()
		case *cellnet.SessionClosed:
			signal.Done(1)
		}
	})

	connector.Start()

	signal.WaitAndExpect("client not closed", 1)
	signal.WaitAndExpect("server not recv closed", 2)

	acc.Stop()
}

func TestUDPServerPositiveClose(t *testing.T) {

	signal := util.NewSignalTester(t)

	acc := peer.NewGenericPeer("udp.Acceptor", "server", udpSes_Address, nil)

	proc.BindProcessorHandler(acc, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionAccepted:
			ev.Session().Close()
			signal.Done(1)
		}
	})

	acc.Start()

	connector := peer.NewGenericPeer("udp.Connector", "client", udpSes_Address, nil)

	proc.BindProcessorHandler(connector, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *cellnet.SessionClosed:
			signal.Done(2)
		}
	})

	connector.Start()

	signal.WaitAndExpect("server not accept", 1)
	signal.WaitAndExpect("client not recv closed closed", 2)

	acc.Stop()
}
