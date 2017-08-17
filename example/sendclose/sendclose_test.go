package sendclose

import (
	"testing"

	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/pb" // 启用pb编码
	"github.com/davyxu/cellnet/proto/binary/coredef"
	"github.com/davyxu/cellnet/proto/pb/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

var accpetor cellnet.Peer

func runServer() {
	queue := cellnet.NewEventQueue()

	accpetor = socket.NewAcceptor(queue).Start("127.0.0.1:7202")
	accpetor.SetName("server")

	cellnet.RegisterMessage(accpetor, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv ", msg.Content)

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

	p := socket.NewConnector(queue).Start("127.0.0.1:7202")
	p.SetName("client.connActiveClose")

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		signal.Done(1)
		log.Debugln("send no close")

		// 连接上发包,告诉服务器不要断开
		ev.Send(&gamedef.TestEchoACK{
			Content: "noclose",
		})

	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())
		signal.Done(2)

		// 客户端主动断开
		ev.Ses.Close()

	})

	cellnet.RegisterMessage(p, "coredef.SessionClosed", func(ev *cellnet.Event) {

		msg := ev.Msg.(*coredef.SessionClosed)

		log.Debugln("close ok!", msg.Result)
		// 正常断开
		signal.Done(3)

	})

	queue.StartLoop()

	signal.WaitAndExpect("TestConnActiveClose not connected", 1)
	signal.WaitAndExpect("TestConnActiveClose not recv msg", 2)
	signal.WaitAndExpect("TestConnActiveClose not close", 3)
}

// 接收封包后被断开
func testRecvDisconnected() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7202")
	p.SetName("client.recvDisconnected")

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		// 连接上发包
		ev.Send(&gamedef.TestEchoACK{
			Content: "RecvDisconnected",
		})

		signal.Done(1)
	})

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(2)

	})

	cellnet.RegisterMessage(p, "coredef.SessionClosed", func(ev *cellnet.Event) {

		// 断开
		signal.Done(3)
	})

	queue.StartLoop()

	signal.WaitAndExpect("TestRecvDisconnected not connected", 1)
	signal.WaitAndExpect("TestRecvDisconnected not recv msg", 2)
	signal.WaitAndExpect("TestRecvDisconnected not closed", 3)

}

func TestConnActiveClose(t *testing.T) {

	signal = util.NewSignalTester(t)

	runServer()

	testConnActiveClose()

	accpetor.Stop()
}

func TestRecvDisconnected(t *testing.T) {

	signal = util.NewSignalTester(t)

	runServer()

	testRecvDisconnected()

	accpetor.Stop()
}
