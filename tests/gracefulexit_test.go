package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"sync"
	"testing"
	"time"
)

const recreateConn_Address = "127.0.0.1:7201"

var recreateConn_Signal *SignalTester

func recreateConn_StartServer() {
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", recreateConn_Address, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *TestEchoACK:

			fmt.Printf("server recv %+v\n", msg)

			ev.Session().Send(&TestEchoACK{
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

	p := peer.NewGenericPeer("tcp.Connector", "client.ConnClose", recreateConn_Address, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			p.Stop()

			time.Sleep(time.Millisecond * 100)

			if times < 3 {
				p.Start()
				times++
			} else {
				recreateConn_Signal.Done(1)
			}
		}
	})

	p.Start()

	queue.StartLoop()

	recreateConn_Signal.WaitAndExpect("not expect times", 1)

	p.Stop()
}

func TestCreateDestroyConnector(t *testing.T) {

	recreateConn_Signal = NewSignalTester(t)

	recreateConn_StartServer()

	runConnClose()
}

const recreateAcc_clientConnection = 3

const recreateAcc_Address = "127.0.0.1:7711"

func TestCreateDestroyAcceptor(t *testing.T) {
	queue := cellnet.NewEventQueue()

	var allAccepted sync.WaitGroup

	p := peer.NewGenericPeer("tcp.Acceptor", "server", recreateAcc_Address, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

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

		p := peer.NewGenericPeer("tcp.Connector", "client.ConnClose", recreateAcc_Address, nil)

		proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		})

		p.Start()

	}

}
