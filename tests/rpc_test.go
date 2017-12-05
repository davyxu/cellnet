package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/tests/proto"
	"github.com/davyxu/cellnet/util"
	"testing"
	"time"
)

const syncRPCAddress = "127.0.0.1:9201"

var syncRPCSignal *util.SignalTester
var asyncRPCSignal *util.SignalTester

var rpcAcceptor cellnet.Peer

func OnRPCServerEvent(raw cellnet.EventParam) cellnet.EventResult {
	switch ev := raw.(type) {
	case rpc.RecvMsgEvent:

		msg := ev.Msg.(*proto.TestEchoACK)

		log.Debugln("server recv rpc ", *msg)

		ev.Reply(&proto.TestEchoACK{
			Msg:   msg.Msg,
			Value: msg.Value,
		})
	}

	return nil
}

func StartRPCServer() {
	queue := cellnet.NewEventQueue()

	rpcAcceptor = cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tcp.Acceptor",
		Queue:       queue,
		PeerAddress: syncRPCAddress,
		PeerName:    "server",
		Event:       OnRPCServerEvent,
	}).Start()

	queue.StartLoop()
}

func OnSyncRPCClientEvent(raw cellnet.EventParam) cellnet.EventResult {
	switch ev := raw.(type) {
	case cellnet.ConnectedEvent:

		for i := 0; i < 2; i++ {

			// 同步阻塞请求必须并发启动，否则客户端无法接收数据
			go func(id int) {

				result, err := rpc.CallSync(ev.Ses, &proto.TestEchoACK{
					Msg:   "sync",
					Value: 1234,
				}, time.Second*5)

				if err != nil {
					syncRPCSignal.Log(err)
					syncRPCSignal.FailNow()
					return
				}

				msg := result.(*proto.TestEchoACK)
				log.Debugln("client sync recv:", msg.Msg, id*100)

				syncRPCSignal.Done(id * 100)

			}(i + 1)
		}

	}

	return nil
}

func OnASyncRPCClientEvent(raw cellnet.EventParam) cellnet.EventResult {
	switch ev := raw.(type) {
	case cellnet.ConnectedEvent:

		for i := 0; i < 2; i++ {

			copy := i + 1

			rpc.Call(ev.Ses, &proto.TestEchoACK{
				Msg:   "async",
				Value: 1234,
			}, time.Second*5, func(feedback interface{}) {

				switch v := feedback.(type) {
				case error:
					asyncRPCSignal.Log(v)
					asyncRPCSignal.FailNow()
				case *proto.TestEchoACK:
					log.Debugln("client sync recv:", v.Msg)
					asyncRPCSignal.Done(copy)
				}

			})

		}

	}

	return nil
}

func StartRPCClient(eventFunc cellnet.EventFunc) {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		PeerType:    "tcp.Connector",
		Queue:       queue,
		PeerAddress: syncRPCAddress,
		PeerName:    "client",
		Event:       eventFunc,
	}).Start()

	queue.StartLoop()
}

func TestSyncRPC(t *testing.T) {

	syncRPCSignal = util.NewSignalTester(t)

	StartRPCServer()

	StartRPCClient(OnSyncRPCClientEvent)
	syncRPCSignal.WaitAndExpect("sync not recv data ", 100, 200)

	rpcAcceptor.Stop()
}

func TestASyncRPC(t *testing.T) {

	asyncRPCSignal = util.NewSignalTester(t)

	StartRPCServer()

	StartRPCClient(OnASyncRPCClientEvent)
	asyncRPCSignal.WaitAndExpect("async not recv data ", 1, 2)

	rpcAcceptor.Stop()
}
