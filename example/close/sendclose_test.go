package sendclose

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

func runServer() {
	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv ")

		// 发包后关闭
		ev.Send(&gamedef.TestEchoACK{
			Content: msg.Content,
		})

		if msg.Content != "noclose" {
			ev.Ses.Close()
		}

	})

	queue.StartLoop()

}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func testConnActiveClose() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(p, "gamedef.SessionConnected", func(ev *cellnet.SessionEvent) {

		signal.Done(1)
		log.Debugln("send no close")

		// 连接上发包,告诉服务器不要断开
		ev.Send(&gamedef.TestEchoACK{
			Content: "noclose",
		})

	})

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())
		signal.Done(2)

		// 客户端主动断开
		ev.Ses.Close()

	})

	socket.RegisterMessage(p, "gamedef.SessionClosed", func(ev *cellnet.SessionEvent) {

		msg := ev.Msg.(*gamedef.SessionClosed)

		log.Debugln("close ok!", msg.Reason)
		// 正常断开
		signal.Done(3)

	})

	queue.StartLoop()

	signal.WaitAndExpect(1, "TestConnActiveClose not connected")
	signal.WaitAndExpect(2, "TestConnActiveClose not recv msg")
	signal.WaitAndExpect(3, "TestConnActiveClose not close")
}

// 接收封包后被断开
func testRecvDisconnected() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(p, "gamedef.SessionConnected", func(ev *cellnet.SessionEvent) {

		// 连接上发包
		ev.Send(&gamedef.TestEchoACK{
			Content: "data",
		})

		signal.Done(1)
	})

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.SessionEvent) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(2)

	})

	socket.RegisterMessage(p, "gamedef.SessionClosed", func(ev *cellnet.SessionEvent) {

		// 断开
		signal.Done(3)
	})

	queue.StartLoop()

	signal.WaitAndExpect(1, "TestRecvDisconnected not connected")
	signal.WaitAndExpect(2, "TestRecvDisconnected not recv msg")
	signal.WaitAndExpect(3, "TestRecvDisconnected not closed")

}

func TestClose(t *testing.T) {

	signal = test.NewSignalTester(t)

	runServer()

	testConnActiveClose()
	testRecvDisconnected()
}
