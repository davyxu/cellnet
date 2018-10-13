package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/rpc"
	"testing"
	"time"
)

const syncRPC_Address = "127.0.0.1:9201"

var (
	syncRPC_Signal  *SignalTester
	asyncRPC_Signal *SignalTester

	rpc_Acceptor cellnet.Peer
)

func rpc_StartServer() {
	queue := cellnet.NewEventQueue()

	rpc_Acceptor = peer.NewGenericPeer("tcp.Acceptor", "server", syncRPC_Address, queue)

	proc.BindProcessorHandler(rpc_Acceptor, "tcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *TestEchoACK:
			log.Debugln("server recv rpc ", *msg)

			ev.(interface {
				Reply(interface{})
			}).Reply(&TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})

		}

	})
	rpc_Acceptor.Start()

	queue.StartLoop()
}
func syncRPC_OnClientEvent(ev cellnet.Event) {

	switch ev.Message().(type) {
	case *cellnet.SessionConnected:
		for i := 0; i < 2; i++ {

			// 同步阻塞请求必须并发启动，否则客户端无法接收数据
			go func(id int) {

				result, err := rpc.CallSync(ev.Session(), &TestEchoACK{
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

func asyncRPC_OnClientEvent(ev cellnet.Event) {

	switch ev.Message().(type) {
	case *cellnet.SessionConnected:
		for i := 0; i < 2; i++ {

			copy := i + 1

			rpc.Call(ev.Session(), &TestEchoACK{
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

func rpc_StartClient(eventFunc func(event cellnet.Event)) {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "client", syncRPC_Address, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", eventFunc)

	p.Start()

	queue.StartLoop()
}

func TestSyncRPC(t *testing.T) {

	syncRPC_Signal = NewSignalTester(t)

	rpc_StartServer()

	rpc_StartClient(syncRPC_OnClientEvent)
	syncRPC_Signal.WaitAndExpect("sync not recv data ", 100, 200)

	rpc_Acceptor.Stop()
}

func TestASyncRPC(t *testing.T) {

	asyncRPC_Signal = NewSignalTester(t)

	rpc_StartServer()

	rpc_StartClient(asyncRPC_OnClientEvent)
	asyncRPC_Signal.WaitAndExpect("async not recv data ", 1, 2)

	rpc_Acceptor.Stop()
}
