package main

import (
	"github.com/davyxu/cellnet"
	jsongamedef "github.com/davyxu/cellnet/proto/json/gamedef" // json逻辑协议
	"github.com/davyxu/cellnet/websocket"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

// 运行服务器, 在浏览器(Chrome)中打开index.html, F12打开调试窗口->Console标签 查看命令行输出
func main() {

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

		ev.Ses.Close()

	})

	queue.StartLoop()

	queue.Wait()

	p.Stop()
}
