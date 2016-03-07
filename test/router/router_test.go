package main

import (
	"testing"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/router"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/test"
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("test")

var signal *test.SignalTester

// 后台服务器
func backendServer() {

	router.DebugMode = true

	pipe := cellnet.NewEventPipe()

	router.StartBackendConnector(pipe, []string{"127.0.0.1:7201"}, "svc->backend")

	router.RegisterSessionMessage("coredef.SessionClosed", func(content interface{}, routerSes cellnet.Session, clientid int64) {
		log.Debugf("client closed router: %d clientid: %d\n", routerSes.ID(), clientid)
	})

	router.RegisterSessionMessage("coredef.TestEchoACK", func(content interface{}, routerSes cellnet.Session, clientid int64) {
		msg := content.(*coredef.TestEchoACK)

		log.Debugf("recv relay,  router: %d clientid: %d\n", routerSes.ID(), clientid)

		signal.Done(2)

		router.SendToClient(routerSes, clientid, &coredef.TestEchoACK{
			Content: msg.Content,
		})
	})

	pipe.Start()
}

// 网关服务器
func routerServer() {

	router.DebugMode = true

	pipe := cellnet.NewEventPipe()

	router.StartBackendAcceptor(pipe, "127.0.0.1:7201", "svc->backend")
	router.StartFrontendAcceptor(pipe, "127.0.0.1:7101", "client->router")

	pipe.Start()

}

// 客户端
func client() {

	pipe := cellnet.NewEventPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7101")

	socket.RegisterSessionMessage(evq, "coredef.SessionConnected", func(content interface{}, ses cellnet.Session) {

		signal.Done(1)

		ack := &coredef.TestEchoACK{
			Content: "hello",
		}
		ses.Send(ack)

		log.Debugf("client send: %s\n", ack.String())

	})

	socket.RegisterSessionMessage(evq, "coredef.TestEchoACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.TestEchoACK)

		log.Debugln("client recv:", msg.String())

		signal.Done(3)
	})

	pipe.Start()

	signal.WaitAndExpect(1, "not connceted to router")
	signal.WaitAndExpect(2, "not recv client msg")
	signal.WaitAndExpect(3, "not recv server msg")

}

func TestGate(t *testing.T) {

	signal = test.NewSignalTester(t)

	routerServer()
	backendServer()
	client()

}
