package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/util"
	"sync"
	"testing"
	"time"
)

const recreateConn_Address = "127.0.0.1:7201"

var recreateConn_Signal *util.SignalTester

func recreateConn_StartServer() {
	queue := cellnet.NewEventQueue()

	p := peer.NewPeer("tcp.Acceptor")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Address", recreateConn_Address)
	pset.SetProperty("Name", "server")
	pset.SetProperty("Queue", queue)

	proc.BindProcessor(p, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *TestEchoACK:

			fmt.Printf("server recv %+v\n", msg)

			ev.BaseSession().(cellnet.Session).Send(&TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})
		}
	})

	p.Start()

	queue.StartLoop()
}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func runConnClose() {
	queue := cellnet.NewEventQueue()

	var times int

	peerIns := peer.NewPeer("tcp.Connector")
	pset := peerIns.(cellnet.PropertySet)
	pset.SetProperty("Address", recreateConn_Address)
	pset.SetProperty("Name", "client.ConnClose")
	pset.SetProperty("Queue", queue)

	proc.BindProcessor(peerIns, "tcp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			peerIns.Stop()

			time.Sleep(time.Millisecond * 100)

			if times < 3 {
				peerIns.Start()
				times++
			} else {
				recreateConn_Signal.Done(1)
			}
		}
	})

	peerIns.Start()

	queue.StartLoop()

	recreateConn_Signal.WaitAndExpect("not expect times", 1)

	peerIns.Stop()
}

func TestCreateDestroyConnector(t *testing.T) {

	recreateConn_Signal = util.NewSignalTester(t)

	recreateConn_StartServer()

	runConnClose()
}

const recreateAcc_clientConnection = 3

const recreateAcc_Address = "127.0.0.1:7711"

func TestCreateDestroyAcceptor(t *testing.T) {
	queue := cellnet.NewEventQueue()

	var allAccepted sync.WaitGroup

	p := peer.NewPeer("tcp.Acceptor")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Address", recreateAcc_Address)
	pset.SetProperty("Name", "server")
	pset.SetProperty("Queue", queue)

	proc.BindProcessor(p, "tcp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionAccepted:

			allAccepted.Done()

		}
	})

	p.Start()

	queue.StartLoop()

	log.Debugln("Start connecting...")
	allAccepted.Add(recreateAcc_clientConnection)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("Close acceptor...")
	p.Stop()

	// 确认所有连接已经断开
	time.Sleep(time.Second)

	log.Debugln("Session count:", p.(cellnet.SessionAccessor).SessionCount())

	p.Start()
	log.Debugln("Start connecting...")
	allAccepted.Add(recreateAcc_clientConnection)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("All done")
}

func runMultiConnection() {

	for i := 0; i < recreateAcc_clientConnection; i++ {

		peerIns := peer.NewPeer("tcp.Connector")
		pset := peerIns.(cellnet.PropertySet)
		pset.SetProperty("Address", recreateAcc_Address)
		pset.SetProperty("Name", "client.ConnClose")

		proc.BindProcessor(peerIns, "tcp.ltv", func(ev cellnet.Event) {

		})

		peerIns.Start()

	}

}
