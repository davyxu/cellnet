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

	evq := cellnet.NewEvQueue()

	socket.NewAcceptor(evq).Start("127.0.0.1:7234")

	socket.RegisterMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String(msg.String()),
		})

	})

}

func client() {

	evq := cellnet.NewEvQueue()

	socket.RegisterMessage(evq, coredef.TestEchoACK{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	socket.RegisterMessage(evq, coredef.ConnectedACK{}, func(ses cellnet.Session, content interface{}) {

		ses.Send(&coredef.TestEchoACK{
			Content: proto.String("hello"),
		})

	})

	socket.NewConnector(evq).Start("127.0.0.1:7234")

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
