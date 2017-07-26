package echowebsocket

import (
	"github.com/davyxu/cellnet"
	jsongamedef "github.com/davyxu/cellnet/proto/json/gamedef" // json逻辑协议
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/cellnet/websocket"
	"github.com/davyxu/golog"
	"testing"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	// 注意, 如果http代理/VPN在运行时可能会导致无法连接, 请关闭

	p := websocket.NewAcceptor(queue).Start("http://127.0.0.1:8801/echo")
	p.SetName("server")

	cellnet.RegisterMessage(p, "coredef.SessionAccepted", func(ev *cellnet.Event) {

		log.Debugln("client accepted")
	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoJsonACK", func(ev *cellnet.Event) {

		msg := ev.Msg.(*jsongamedef.TestEchoJsonACK)

		log.Debugln(msg.Content)

		ev.Send(&jsongamedef.TestEchoJsonACK{Content: "roger"})
	})

	queue.StartLoop()
}

func client() {

	queue := cellnet.NewEventQueue()

	p := websocket.NewConnector(queue).Start("ws://127.0.0.1:8801/echo")
	p.SetName("client")

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		// 发送消息, 底层自动选择pb编码
		ev.Send(&jsongamedef.TestEchoJsonACK{
			Content: "hello",
		})

		signal.Done(1)

	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoJsonACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*jsongamedef.TestEchoJsonACK)

		log.Debugln("client recv:", msg.Content)

		signal.Done(2)
	})

	queue.StartLoop()

	signal.WaitAndExpect("not recv data", 1, 2)

}

func TestWebsocketEcho(t *testing.T) {

	signal = util.NewSignalTester(t)

	server()

	client()

}
