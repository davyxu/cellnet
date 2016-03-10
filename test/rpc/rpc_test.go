package rpc

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
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
	p.Start("127.0.0.1:7201")
	rpc.InstallServer(p)

	rpc.RegisterMessage(p, "coredef.TestEchoACK", func(resp rpc.Response, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		resp.Feedback(&coredef.TestEchoACK{
			Content: msg.String(),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewConnector(pipe)
	p.SetName("client")
	p.Start("127.0.0.1:7201")

	rpc.InstallClient(p)

	socket.RegisterSessionMessage(p, "coredef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		rpc.Call(p, &coredef.TestEchoACK{
			Content: "rpc hello",
		}, func(msg *coredef.TestEchoACK) {

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
