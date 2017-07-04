package echo_pb

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/pb/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *util.SignalTester

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

func server() {

	queue := cellnet.NewEventQueue()

	peer := socket.NewAcceptor(queue).Start("127.0.0.1:7301")
	_, send := peer.HandlerList()

	// 在原有流程中, 插入固定消息回调
	recvList := socket.BuildRecvHandler(new(RecvMessageHandler), peer)

	peer.SetHandlerList(recvList, send)

	cellnet.RegisterMessage(peer, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("server recv:", msg.Content)

		signal.Done(1)
	})

	queue.StartLoop()

}

// 客户端为了逻辑编写方便, 依然使用dispatcher配合闭包消息处理函数方式
func client() {

	queue := cellnet.NewEventQueue()

	dh := socket.NewConnector(queue).Start("127.0.0.1:7301")

	cellnet.RegisterMessage(dh, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		log.Debugln("client recv:", msg.Content)

		signal.Done(1)
	})

	cellnet.RegisterMessage(dh, "coredef.SessionConnected", func(ev *cellnet.Event) {

		log.Debugln("client connected")

		// 发送消息, 底层自动选择pb编码
		ev.Send(&gamedef.TestEchoACK{
			Content: "hello",
		})

	})

	queue.StartLoop()

	signal.WaitAndExpect("not recv data", 1)

}

func TestEcho(t *testing.T) {

	signal = util.NewSignalTester(t)

	server()

	client()

}
