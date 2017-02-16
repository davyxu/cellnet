package sproto_echo

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec/sproto"
	"github.com/davyxu/cellnet/example"
	"github.com/davyxu/cellnet/proto/sproto/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	evd := socket.NewAcceptor(queue).Start("127.0.0.1:7201")
	evd.SetName("server")

	cellnet.RegisterMessage(evd, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
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

	dh := socket.NewConnector(queue).Start("127.0.0.1:7201")
	dh.SetName("client")

	cellnet.RegisterMessage(dh, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)

		signal.Done(1)
	})

	cellnet.RegisterMessage(dh, "gamedef.SessionConnected", func(ev *cellnet.SessionEvent) {

		log.Debugln("client connected")

		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	cellnet.RegisterMessage(dh, "gamedef.SessionConnectFailed", func(ev *cellnet.SessionEvent) {

		msg := ev.Msg.(*gamedef.SessionConnectFailed)

		log.Debugln(msg.Reason)

	})

	queue.StartLoop()

	signal.WaitAndExpect(1, "not recv data")

}

func TestSprotoEcho(t *testing.T) {

	socket.DefaultCodec = "sproto"
	sprotocodec.AutoRegisterMessageMeta(gamedef.SProtoStructs)

	signal = test.NewSignalTester(t)

	server()

	client()

}
