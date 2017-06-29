package echo_sproto

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/binary/coredef"
	"github.com/davyxu/cellnet/proto/sproto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewAcceptor(queue).Start("127.0.0.1:7401")
	evd.SetName("server")

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.Content)

		ev.Send(&gamedef.TestEchoACK{
			Content: msg.Content,
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	dh := socket.NewConnector(queue).Start("127.0.0.1:7401")
	dh.SetName("client")

	cellnet.RegisterMessage(dh, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)

		signal.Done(1)
	})

	cellnet.RegisterMessage(dh, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	cellnet.RegisterMessage(dh, "coredef.SessionConnectFailed", func(ev *cellnet.Event) {

		msg := ev.Msg.(*coredef.SessionConnectFailed)

		log.Debugln(msg.Result)

	})

	queue.StartLoop()

	signal.WaitAndExpect("not recv data", 1)

}

func TestSprotoEcho(t *testing.T) {

	signal = util.NewSignalTester(t)

	server()

	client()

}
