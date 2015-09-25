package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

var done = make(chan bool)

func server() {

	pipe := cellnet.NewEvPipe()

	evq := socket.NewAcceptor(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.String()),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEvPipe()

	evq := socket.NewConnector(pipe).Start("127.0.0.1:7234")

	socket.RegisterSessionMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	socket.RegisterSessionMessage(evq, coredef.SessionConnected{}, func(ses cellnet.Session, content interface{}) {

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String("hello"),
		})

	})

	pipe.Start()
}

func main() {

	server()

	client()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		log.Println("time out")
	}

}
