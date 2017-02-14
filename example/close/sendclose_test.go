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

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		// 发包后关闭
		ses.Send(&gamedef.TestEchoACK{
			Content: msg.Content,
		})

		if msg.Content != "noclose" {
			ses.Close()
		}

	})

	queue.StartLoop()

}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func testConnActiveClose() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7201")

	socket.RegisterMessage(p, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		signal.Done(1)
		// 连接上发包,告诉服务器不要断开
		ses.Send(&gamedef.TestEchoACK{
			Content: "noclose",
		})

	})

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())
		signal.Done(2)

		// 客户端主动断开
		ses.Close()

	})

	socket.RegisterMessage(p, "gamedef.SessionClosed", func(content interface{}, ses cellnet.Session) {

		msg := content.(*gamedef.SessionClosed)

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

	socket.RegisterMessage(p, "gamedef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		// 连接上发包
		ses.Send(&gamedef.TestEchoACK{
			Content: "data",
		})

		signal.Done(1)
	})

	socket.RegisterMessage(p, "gamedef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(2)

	})

	socket.RegisterMessage(p, "gamedef.SessionClosed", func(content interface{}, ses cellnet.Session) {

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
