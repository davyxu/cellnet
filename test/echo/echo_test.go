package echo

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

func server() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewAcceptor(pipe).Start("127.0.0.1:7201")

	socket.RegisterSessionMessage(evq, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.String())

		ses.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7201")

	socket.RegisterSessionMessage(evq, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(1)
	})

	socket.RegisterSessionMessage(evq, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		ses.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	pipe.Start()

	signal.WaitAndExpect(1, "not recv data")

}

func TestEcho(t *testing.T) {

	signal = test.NewSignalTester(t)

	server()

	client()

}
