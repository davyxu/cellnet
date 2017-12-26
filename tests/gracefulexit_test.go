package tests

import (
	"testing"

	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/util"
	"sync"
	"time"
)

const createDestoryConnectorAddress = "127.0.0.1:7201"

var singalConnector *util.SignalTester

func StartCreateDestoryServer() {
	queue := cellnet.NewEventQueue()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    createDestoryConnectorAddress,
		PeerName:       "server",
		InboundEvent: func(raw cellnet.EventParam) cellnet.EventResult {

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

	queue.StartLoop()
}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func runConnClose() {
	queue := cellnet.NewEventQueue()

	var times int

	var peer cellnet.Peer
	peer = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Connector",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    createDestoryConnectorAddress,
		PeerName:       "client.ConnClose",
		InboundEvent: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionConnected:
					peer.Stop()

					time.Sleep(time.Millisecond * 500)

					if times < 3 {
						peer.Start()
						times++
					} else {
						singalConnector.Done(1)
					}
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	singalConnector.WaitAndExpect("not expect times", 1)

	peer.Stop()
}

func TestCreateDestroyConnector(t *testing.T) {

	singalConnector = util.NewSignalTester(t)
	singalConnector.SetTimeout(time.Second * 3)

	StartCreateDestoryServer()

	runConnClose()
}

const clientConnectionCount = 3

const createDestoryAcceptorAddress = "127.0.0.1:7711"

func TestCreateDestroyAcceptor(t *testing.T) {
	queue := cellnet.NewEventQueue()

	var allAccepted sync.WaitGroup
	p := cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    createDestoryAcceptorAddress,
		PeerName:       "server",
		InboundEvent: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionAccepted:

					allAccepted.Done()

				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	log.Debugln("Start connecting...")
	allAccepted.Add(clientConnectionCount)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("Close acceptor...")
	p.Stop()

	// 确认所有连接已经断开
	time.Sleep(time.Second)

	log.Debugln("Session count:", p.SessionCount())

	p.Start()
	log.Debugln("Start connecting...")
	allAccepted.Add(clientConnectionCount)
	runMultiConnection()

	log.Debugln("Wait all accept...")
	allAccepted.Wait()

	log.Debugln("All done")
}

func runMultiConnection() {

	for i := 0; i < clientConnectionCount; i++ {

		cellnet.CreatePeer(cellnet.PeerConfig{
			PeerType:       "tcp.Connector",
			EventProcessor: "tcp.ltv",
			PeerAddress:    createDestoryAcceptorAddress,
			PeerName:       "client.ConnClose",
		}).Start()
	}

}
