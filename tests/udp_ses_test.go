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

	acceptor := peer.NewPeer("udp.Acceptor")
	pset := acceptor.(cellnet.PropertySet)
	pset.SetProperty("Address", udpSes_Address)
	pset.SetProperty("Name", "server")

	proc.BindProcessor(acceptor, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:
			signal.Done(2)
		}
	})

	acceptor.Start()

	connector := peer.NewPeer("udp.Connector")
	pset2 := connector.(cellnet.PropertySet)
	pset2.SetProperty("Address", udpSes_Address)
	pset2.SetProperty("Name", "client")

	proc.BindProcessor(connector, "udp.ltv", func(ev cellnet.Event) {

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

	acceptor.Stop()
}

func TestUDPServerPositiveClose(t *testing.T) {

	signal := util.NewSignalTester(t)

	acceptor := peer.NewPeer("udp.Acceptor")
	pset := acceptor.(cellnet.PropertySet)
	pset.SetProperty("Address", udpSes_Address)
	pset.SetProperty("Name", "server")

	proc.BindProcessor(acceptor, "udp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionAccepted:
			ev.Session().Close()
			signal.Done(1)
		}
	})

	acceptor.Start()

	connector := peer.NewPeer("udp.Connector")
	pset2 := connector.(cellnet.PropertySet)
	pset2.SetProperty("Address", udpSes_Address)
	pset2.SetProperty("Name", "client")

	proc.BindProcessor(connector, "udp.ltv", func(ev cellnet.Event) {

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

	acceptor.Stop()
}
