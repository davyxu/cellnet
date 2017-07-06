package gracefulexit

import (
	"testing"

	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/pb" // 启用pb编码
	"github.com/davyxu/cellnet/proto/pb/gamedef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/util"
	"time"
)

var singalConnector *util.SignalTester

func runEchoServer() {
	queue := cellnet.NewEventQueue()

	p := socket.NewAcceptor(queue).Start("127.0.0.1:7201")
	p.SetName("server")

	cellnet.RegisterMessage(p, "gamedef.TestEchoACK", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.TestEchoACK)

		// 发包后关闭
		ev.Send(&gamedef.TestEchoACK{
			Content: msg.Content,
		})

	})

	queue.StartLoop()

}

// 客户端连接上后, 主动断开连接, 确保连接正常关闭
func runConnClose() {

	queue := cellnet.NewEventQueue()

	p := socket.NewConnector(queue).Start("127.0.0.1:7201")
	p.SetName("client.ConnClose")

	var times int

	cellnet.RegisterMessage(p, "coredef.SessionConnected", func(ev *cellnet.Event) {

		p.Stop()

		time.Sleep(time.Millisecond * 500)

		if times < 3 {
			p.Start(p.Address())
			times++
		} else {
			singalConnector.Done(1)
		}

	})

	queue.StartLoop()

	singalConnector.WaitAndExpect("not expect times", 1)

	p.Stop()
}

func TestCreateDestroyConnector(t *testing.T) {

	singalConnector = util.NewSignalTester(t)
	singalConnector.SetTimeout(time.Second * 3)

	runEchoServer()

	runConnClose()
}
