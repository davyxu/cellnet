package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/gate"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"

	"runtime"
)

var done = make(chan bool)

func backendServer() {

	gate.DebugMode = true

	pipe := cellnet.NewEvPipe()

	gate.StartGateConnector(pipe, []string{"127.0.0.1:7201"})

	gate.RegisterSessionMessage(coredef.SessionClosed{}, func(gateSes cellnet.Session, clientid int64, content interface{}) {
		log.Printf("client closed gate: %d clientid: %d\n", gateSes.ID(), clientid)
	})

	gate.RegisterSessionMessage(coredef.TestEchoACK{}, func(gateSes cellnet.Session, clientid int64, content interface{}) {
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

	//	fmt.Println(socket.Event_SessionAccepted, "socket.Event_SessionAccepted")
	//	fmt.Println(socket.Event_SessionClosed, "socket.Event_SessionClosed")
	//	fmt.Println(socket.Event_SessionConnected, "socket.Event_SessionConnected")

	gate.DebugMode = true

	pipe := cellnet.NewEvPipe()

	gate.StartBackendAcceptor(pipe, "127.0.0.1:7201")
	gate.StartClientAcceptor(pipe, "127.0.0.1:7101")

	socket.RegisterSessionMessage(gate.ClientAcceptor, coredef.SessionAccepted{}, func(ses cellnet.Session, content interface{}) {

		log.Println("client accepted", ses.ID())

	})

	pipe.Start()

	<-done
}

func client() {

	pipe := cellnet.NewEvPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7101")

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(ses cellnet.Session, content interface{}) {

		ack := &coredef.TestEchoACK{
			Content: proto.String("hello"),
		}
		ses.Send(ack)

		log.Printf("client send: %s\n", ack.String())

	})

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	pipe.Start()

	<-done
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	mode := flag.String("mode", "gate", "specify the mode of this test")

	flag.Parse()

	switch *mode {
	case "gate":
		gateServer()

	case "client":
		client()
	case "backend":
		backendServer()
	}

}
