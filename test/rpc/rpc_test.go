package rpc

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

func server() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewAcceptor(pipe)
	p.SetName("server")
	p.Start("127.0.0.1:9201")

	rpc.RegisterMessage(p, "gamedef.TestEchoACK", func(resp rpc.Response, content interface{}) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		resp.Feedback(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewConnector(pipe)
	p.SetName("client")
	p.Start("127.0.0.1:9201")

	socket.RegisterSessionMessage(p, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		rpc.Call(p, &gamedef.TestEchoACK{
			Content: "rpc async call",
		}, func(msg *gamedef.TestEchoACK) {

			log.Debugln("client recv", msg.Content)

			signal.Done(1)
		})

	})

	pipe.Start()

	signal.WaitAndExpect(1, "not recv data")
}

func TestRPC(t *testing.T) {

	signal = test.NewSignalTester(t)

	server()

	client()

}
