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

const recreateConn_Address = "127.0.0.1:7201"

var recreateConn_Signal *util.SignalTester

func recreateConn_StartServer() {
	queue := cellnet.NewEventQueue()

	cellnet.CreatePeer(cellnet.CommunicatePeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    recreateConn_Address,
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

	queue.StartLoop()
}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func runConnClose() {
	queue := cellnet.NewEventQueue()

	var times int

	var peer cellnet.Peer
	peer = cellnet.CreatePeer(cellnet.CommunicatePeerConfig{
		PeerType:       "tcp.Connector",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    recreateConn_Address,
		PeerName:       "client.ConnClose",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*cellnet.RecvMsgEvent)
			if ok {
				switch ev.Msg.(type) {
				case *comm.SessionConnected:
					peer.Stop()

					time.Sleep(time.Millisecond * 100)

					if times < 3 {
						peer.Start()
						times++
					} else {
						recreateConn_Signal.Done(1)
					}
				}
			}

			return nil
		},
	}).Start()

	queue.StartLoop()

	recreateConn_Signal.WaitAndExpect("not expect times", 1)

	peer.Stop()
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
	p := cellnet.CreatePeer(cellnet.CommunicatePeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    recreateAcc_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

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

		cellnet.CreatePeer(cellnet.CommunicatePeerConfig{
			PeerType:       "tcp.Connector",
			EventProcessor: "tcp.ltv",
			PeerAddress:    recreateAcc_Address,
			PeerName:       "client.ConnClose",
		}).Start()
	}

}
