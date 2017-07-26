package echopb

import (
	"testing"

	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/pb"                     // 启用pb编码
	"github.com/davyxu/cellnet/proto/binary/coredef"           // 底层系统事件
	jsongamedef "github.com/davyxu/cellnet/proto/json/gamedef" // json逻辑协议
	"github.com/davyxu/cellnet/proto/pb/gamedef"               // pb逻辑协议
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

func server() {

	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue).Start("127.0.0.1:7701")
	p.SetName("server")

	// 混合协议支持, 接收pb编码的消息
	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.Content)

		ev.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	// 混合协议支持, 接收json编码的消息
	cellnet.RegisterMessage(p, "gamedef.TestEchoJsonACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*jsongamedef.TestEchoJsonACK)

		log.Debugln("server recv json:", msg.Content)

		ev.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})

	})

	queue.StartLoop()

}

func client() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7701")
	p.SetName("client")

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)

		signal.Done(1)
	})

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		// 发送消息, 底层自动选择pb编码
		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

		// 发送消息, 底层自动选择json编码
		ev.Send(&jsongamedef.TestEchoJsonACK{
			Content: "hello json",
		})

	})

	cellnet.RegisterMessage(p, "coredef.SessionConnectFailed", func(ev *cellnet.Event) {

		msg := ev.Msg.(*coredef.SessionConnectFailed)

		log.Debugln(msg.Result)

	})

	queue.StartLoop()

	signal.WaitAndExpect("not recv data", 1)

}

func TestEcho(t *testing.T) {

	signal = util.NewSignalTester(t)

	server()

	client()

}
