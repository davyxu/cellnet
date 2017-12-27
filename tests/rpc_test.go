package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/comm"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/util"
	"testing"
	"time"
)

const syncRPCAddress = "127.0.0.1:9201"

var syncRPCSignal *util.SignalTester
var asyncRPCSignal *util.SignalTester

var rpcAcceptor cellnet.Peer

func StartRPCServer() {
	queue := cellnet.NewEventQueue()

	rpcAcceptor = cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:       "tcp.Acceptor",
		EventProcessor: "tcp.ltv",
		Queue:          queue,
		PeerAddress:    syncRPCAddress,
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

func onSyncRPCClientEvent(raw cellnet.EventParam) cellnet.EventResult {

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
						syncRPCSignal.Log(err)
						syncRPCSignal.FailNow()
						return
					}

					msg := result.(*TestEchoACK)
					log.Debugln("client sync recv:", msg.Msg, id*100)

					syncRPCSignal.Done(id * 100)

				}(i + 1)
			}
		}
	}

	return nil
}

func onASyncRPCClientEvent(raw cellnet.EventParam) cellnet.EventResult {

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
						asyncRPCSignal.Log(v)
						asyncRPCSignal.FailNow()
					case *TestEchoACK:
						log.Debugln("client sync recv:", v.Msg)
						asyncRPCSignal.Done(copy)
					}

				})

			}
		}
	}

	return nil
}

func StartRPCClient(eventFunc cellnet.EventProc) {
	queue := cellnet.NewEventQueue()

	cellnet.CreatePeer(cellnet.PeerConfig{
		PeerType:        "tcp.Connector",
		EventProcessor:  "tcp.ltv",
		Queue:           queue,
		PeerAddress:     syncRPCAddress,
		PeerName:        "client",
		UserInboundProc: eventFunc,
	}).Start()

	queue.StartLoop()
}

func TestSyncRPC(t *testing.T) {

	syncRPCSignal = util.NewSignalTester(t)

	StartRPCServer()

	StartRPCClient(onSyncRPCClientEvent)
	syncRPCSignal.WaitAndExpect("sync not recv data ", 100, 200)

	rpcAcceptor.Stop()
}

func TestASyncRPC(t *testing.T) {

	asyncRPCSignal = util.NewSignalTester(t)

	StartRPCServer()

	StartRPCClient(onASyncRPCClientEvent)
	asyncRPCSignal.WaitAndExpect("async not recv data ", 1, 2)

	rpcAcceptor.Stop()
}
