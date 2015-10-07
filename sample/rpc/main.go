package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"github.com/davyxu/cellnet/rpc"
	"log"
	"time"
)

var done = make(chan bool)

func server() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewAcceptor(pipe).Start("127.0.0.1:7234")
	rpc.InstallServer(p)

	rpc.RegisterMessage(p, coredef.TestEchoACK{}, func(resp rpc.Response, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		resp.Feedback(&coredef.TestEchoACK{
			Content: proto.String(msg.String()),
		})

	})

	pipe.Start()

}

func client() {

	pipe := cellnet.NewEventPipe()

	p := socket.NewConnector(pipe).Start("127.0.0.1:7234")
	
	rpc.InstallClient(p)

	socket.RegisterSessionMessage(p, coredef.SessionConnected{}, func(ses cellnet.Session, content interface{}) {
		
		
		rpc.Call( p, &coredef.TestEchoACK{
			Content: proto.String("rpc hello"),
		},  func( msg *coredef.TestEchoACK ){

			
			log.Println("client recv", msg.GetContent())
			
			done <- true
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
