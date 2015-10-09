package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/gate"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/golang/protobuf/proto"
	"log"
	"testing"
)

var signal *test.SignalTester

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

		signal.Done(2)

		gate.SendToClient(gateSes, clientid, &coredef.TestEchoACK{
			Content: proto.String(msg.GetContent()),
		})
	})

	pipe.Start()
}

// 网关服务器
func gateServer() {

	gate.DebugMode = true

	pipe := cellnet.NewEventPipe()

	gate.StartBackendAcceptor(pipe, "127.0.0.1:7201")
	gate.StartClientAcceptor(pipe, "127.0.0.1:7101")

	pipe.Start()

}

// 客户端
func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7101")

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(content interface{}, ses cellnet.Session) {

		signal.Done(1)

		ack := &coredef.TestEchoACK{
			Content: proto.String("hello"),
		}
		ses.Send(ack)

		log.Printf("client send: %s\n", ack.String())

	})

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		signal.Done(3)
	})

	pipe.Start()

	signal.WaitAndExpect(1, "not connceted to gate")
	signal.WaitAndExpect(2, "not recv client msg")
	signal.WaitAndExpect(3, "not recv server msg")

}

func TestGate(t *testing.T) {

	signal = test.NewSignalTester(t)

	gateServer()
	backendServer()
	client()

}
