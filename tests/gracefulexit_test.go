package tests

import (
	"testing"

	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm/tcppkt"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"sync"
	"time"
)

const createDestoryConnectorAddress = "127.0.0.1:7201"

var singalConnector *util.SignalTester

func StartCreateDestoryServer() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Acceptor",
		Queue:       queue,
		PeerAddress: createDestoryConnectorAddress,
		PeerName:    "server",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {
				case *proto.TestEchoACK:

					fmt.Printf("server recv %+v\n", msg)

					ev.Ses.Send(&proto.TestEchoACK{
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
	peer = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Connector",
		Queue:       queue,
		PeerAddress: createDestoryConnectorAddress,
		PeerName:    "client.ConnClose",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *tcppkt.SessionConnected:
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

const createDestoryAcceptorAddress = "127.0.0.1:7701"

func TestCreateDestroyAcceptor(t *testing.T) {
	queue := cellnet.NewEventQueue()

	var allAccepted sync.WaitGroup
	p := cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "ltv.tcp.Acceptor",
		Queue:       queue,
		PeerAddress: createDestoryAcceptorAddress,
		PeerName:    "server",
		Event: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *tcppkt.SessionAccepted:

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

		cellnet.NewPeer(cellnet.PeerConfig{
			PeerType:    "ltv.tcp.Connector",
			PeerAddress: createDestoryAcceptorAddress,
			PeerName:    "client.ConnClose",
		}).Start()
	}

}
