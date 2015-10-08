package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/gate"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"runtime"
)

var done = make(chan bool)

// 后台服务器
func backendServer() {

	gate.DebugMode = true

	pipe := cellnet.NewEventPipe()

	gate.StartGateConnector(pipe, []string{"127.0.0.1:7201"})

	gate.RegisterSessionMessage(coredef.SessionClosed{}, func(content interface{}, gateSes cellnet.Session, clientid int64) {
		log.Printf("client closed gate: %d clientid: %d\n", gateSes.ID(), clientid)
	})

	gate.RegisterSessionMessage(coredef.TestEchoACK{}, func(content interface{}, gateSes cellnet.Session, clientid int64) {
		msg := content.(*coredef.TestEchoACK)

		log.Printf("recv relay,  gate: %d clientid: %d\n", gateSes.ID(), clientid)

		gate.SendToClient(gateSes, clientid, &coredef.TestEchoACK{
			Content: proto.String(msg.GetContent()),
		})
	})

	pipe.Start()

	<-done
}

// 网关服务器
func gateServer() {

	gate.DebugMode = true

	pipe := cellnet.NewEventPipe()

	gate.StartBackendAcceptor(pipe, "127.0.0.1:7201")
	gate.StartClientAcceptor(pipe, "127.0.0.1:7101")

	socket.RegisterSessionMessage(gate.ClientAcceptor, coredef.SessionAccepted{}, func(content interface{}, ses cellnet.Session) {

		log.Println("client accepted", ses.ID())

	})

	pipe.Start()

	<-done
}

// 客户端
func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7101")

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(content interface{}, ses cellnet.Session) {

		ack := &coredef.TestEchoACK{
			Content: proto.String("hello"),
		}
		ses.Send(ack)

		log.Printf("client send: %s\n", ack.String())

	})

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	pipe.Start()

	<-done
}

// 启动顺序:
// 网关服务器: gate gate
// 后台服务器: gate backend
// 客户端: gate client
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) <= 1 {
		return
	}

	switch os.Args[1] {
	case "gate":
		gateServer()
	case "client":
		client()
	case "backend":
		backendServer()
	}

}
