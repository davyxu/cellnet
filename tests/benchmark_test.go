package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/util"
	"testing"
)

func server() {
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "server", "127.0.0.1:7701", queue)

	proc.BindProcessor(p, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *TestEchoACK:

			ev.Session().Send(&TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			})
		}

	})

	p.Start()

	queue.StartLoop()
}

func client() {

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:7701", queue)

	proc.BindProcessor(p, "tcp.ltv", func(event cellnet.Event) {})

	p.Start()

	queue.StartLoop()

	util.SyncRecvEvent(p, func(ev cellnet.Event) {
		msg := ev.Message().(*cellnet.SessionConnected)
		msg = msg

		ev.Session().Send(&TestEchoACK{
			Msg:   "hello",
			Value: 1234,
		})

	})

	util.SyncRecvEvent(p, func(ev cellnet.Event) {

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
