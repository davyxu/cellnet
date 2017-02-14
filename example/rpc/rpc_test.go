package rpc

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/example"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue)
	p.SetName("server")
	p.Start("127.0.0.1:9201")

	rpc.RegisterMessage(p, "gamedef.TestEchoACK", func(content interface{}, resp rpc.Response) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		resp.Feedback(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue)
	p.SetName("client")
	p.Start("127.0.0.1:9201")

	socket.RegisterMessage(p, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		rpc.Call(p, &gamedef.TestEchoACK{
			Content: "rpc async call",
		}, func(msg *gamedef.TestEchoACK) {

			log.Debugln("client recv", msg.Content)

			signal.Done(1)
		})

	})

	queue.StartLoop()

	signal.WaitAndExpect(1, "not recv data")
}

func TestRPC(t *testing.T) {

	signal = test.NewSignalTester(t)

	server()

	client()

}
