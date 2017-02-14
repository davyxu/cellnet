package echo

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/example"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewAcceptor(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		ses.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewConnector(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(evd, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(1)
	})

	socket.RegisterMessage(evd, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		ses.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	socket.RegisterMessage(evd, "gamedef.SessionConnectFailed", func(content interface{}, ses cellnet.Session) {

		msg := content.(*gamedef.SessionConnectFailed)

		log.Debugln(msg.Reason)

	})

	queue.StartLoop()

	signal.WaitAndExpect(1, "not recv data")

}

func TestEcho(t *testing.T) {

	signal = test.NewSignalTester(t)

	server()

	client()

}
