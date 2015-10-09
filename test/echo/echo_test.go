package echo

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/golang/protobuf/proto"
	"log"
	"testing"
)

var signal *test.SignalTester

func server() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewAcceptor(pipe).Start("127.0.0.1:7201")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.String()),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7201")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		signal.Done(1)
	})

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(content interface{}, ses cellnet.Session) {

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String("hello"),
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
