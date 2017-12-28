package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/comm/rpc"
	"github.com/davyxu/cellnet/util"
	"testing"
	"time"
)

const syncRPC_Address = "127.0.0.1:9201"

var syncRPC_Signal *util.SignalTester
var asyncRPC_Signal *util.SignalTester

var rpc_Acceptor cellnet.Peer

func rpc_StartServer() {
	queue := cellnet.NewEventQueue()

	rpc_Acceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    syncRPC_Address,
		PeerName:       "server",
		UserInboundProc: func(raw cellnet.EventParam) cellnet.EventResult {

			ev, ok := raw.(*rpc.RecvMsgEvent)
			if ok {
				switch msg := ev.Msg.(type) {

				case *TestEchoACK:
					log.Debugln("server recv rpc ", *msg)

					ev.Reply(&TestEchoACK{
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

func syncRPC_OnClientEvent(raw cellnet.EventParam) cellnet.EventResult {

	ev, ok := raw.(*cellnet.RecvMsgEvent)
	if ok {
		switch ev.Msg.(type) {
		case *comm.SessionConnected:
			for i := 0; i < 2; i++ {

				// 同步阻塞请求必须并发启动，否则客户端无法接收数据
				go func(id int) {

					result, err := rpc.CallSync(ev.Ses, &TestEchoACK{
						Msg:   "sync",
						Value: 1234,
					}, time.Second*5)

					if err != nil {
						syncRPC_Signal.Log(err)
						syncRPC_Signal.FailNow()
						return
					}

					msg := result.(*TestEchoACK)
					log.Debugln("client sync recv:", msg.Msg, id*100)

					syncRPC_Signal.Done(id * 100)

				}(i + 1)
			}
		}
	}

	return nil
}

func asyncRPC_OnClientEvent(raw cellnet.EventParam) cellnet.EventResult {

	ev, ok := raw.(*cellnet.RecvMsgEvent)
	if ok {
		switch ev.Msg.(type) {
		case *comm.SessionConnected:
			for i := 0; i < 2; i++ {

				copy := i + 1

				rpc.Call(ev.Ses, &TestEchoACK{
					Msg:   "async",
					Value: 1234,
				}, time.Second*5, func(feedback interface{}) {

					switch v := feedback.(type) {
					case error:
						asyncRPC_Signal.Log(v)
						asyncRPC_Signal.FailNow()
					case *TestEchoACK:
						log.Debugln("client sync recv:", v.Msg)
						asyncRPC_Signal.Done(copy)
					}

				})

			}
		}
	}

	return nil
}

func rpc_StartClient(eventFunc cellnet.EventProc) {
	queue := cellnet.NewEventQueue()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:        "tcp.Connector",
		EventProcessor:  "tcp.ltv",
		Queue:           queue,
		PeerAddress:     syncRPC_Address,
		PeerName:        "client",
		UserInboundProc: eventFunc,
	}).Start()

	queue.StartLoop()
}

func TestSyncRPC(t *testing.T) {

	syncRPC_Signal = util.NewSignalTester(t)

	rpc_StartServer()

	rpc_StartClient(syncRPC_OnClientEvent)
	syncRPC_Signal.WaitAndExpect("sync not recv data ", 100, 200)

	rpc_Acceptor.Stop()
}

func TestASyncRPC(t *testing.T) {

	asyncRPC_Signal = util.NewSignalTester(t)

	rpc_StartServer()

	rpc_StartClient(asyncRPC_OnClientEvent)
	asyncRPC_Signal.WaitAndExpect("async not recv data ", 1, 2)

	rpc_Acceptor.Stop()
}
