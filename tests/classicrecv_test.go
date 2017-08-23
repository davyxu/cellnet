package tests

import (
	"testing"

	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/pb" // 启用pb编码
	"github.com/davyxu/cellnet/proto/pb/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
)

var classicalRecvSignal *util.SignalTester

type RecvMessageHandler struct {
}

func (self *RecvMessageHandler) Call(ev *cellnet.Event) {

	onServerMessage(ev)

}

// 接收peer的所有消息, 使用这种传统的结构可以方便做服务器热更新
func onServerMessage(ev *cellnet.Event) {

	switch msg := ev.Msg.(type) {
	case *gamedef.TestEchoACK:

		log.Debugln("classic server recv:", msg.Content)

		ev.Send(&gamedef.TestEchoACK{
			Content: msg.String(),
		})
	}

}

var classicalAcceptor cellnet.Peer

func classicalServer() {

	queue := cellnet.NewEventQueue()

	classicalAcceptor = socket.NewAcceptor(queue).Start("127.0.0.1:7701")
	classicalAcceptor.SetName("server")

	// 添加一条新的处理链
	classicalAcceptor.AddChainRecv(cellnet.NewHandlerChain(
		cellnet.StaticDecodePacketHandler(),
		cellnet.NewQueuePostHandler(queue, new(RecvMessageHandler)),
	))

	cellnet.RegisterMessage(classicalAcceptor, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.Content)

		classicalRecvSignal.Done(1)
	})

	queue.StartLoop()

}

// 客户端为了逻辑编写方便, 依然使用dispatcher配合闭包消息处理函数方式
func classicalClient() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7701")
	p.SetName("client")

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)

		classicalRecvSignal.Done(1)
	})

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		// 发送消息, 底层自动选择pb编码
		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	queue.StartLoop()

	classicalRecvSignal.WaitAndExpect("not recv data", 1)

}

func TestClassicalRecv(t *testing.T) {

	classicalRecvSignal = util.NewSignalTester(t)

	classicalServer()

	classicalClient()

	classicalAcceptor.Stop()
}
