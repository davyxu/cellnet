package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"testing"
)

func server() {
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:7701", queue)

	proc.BindProcessorHandler(p, "tcp.ltv", proc.GlobalDispatcher.OnEvent)

	proc.RegisterMessage(p, "tests.TestEchoACK", func(ev cellnet.Event) {

		msg := ev.Message().(*TestEchoACK)

		ev.Session().Send(&TestEchoACK{
			Msg:   msg.Msg,
			Value: msg.Value,
		})
	})

	p.Start()

	queue.StartLoop()
}

func client() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:7701", queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(event cellnet.Event) {})

	p.Start()

	queue.StartLoop()

	rv := proc.NewSyncReceiver(p)

	rv.Recv(func(ev cellnet.Event) {
		msg := ev.Message().(*cellnet.SessionConnected)
		msg = msg

		ev.Session().Send(&TestEchoACK{
			Msg:   "hello",
			Value: 1234,
		})

	}).Recv(func(ev cellnet.Event) {

		msg := ev.Message().(*TestEchoACK)

		println(msg)

	})

}

func BenchmarkEcho(b *testing.B) {
	b.ReportAllocs()

	server()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client()
	}

}
