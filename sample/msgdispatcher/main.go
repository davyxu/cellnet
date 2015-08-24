package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

var done = make(chan bool)

func server() {

	disp := dispatcher.NewPacketDispatcher()

	dispatcher.RegisterMessage(disp, coredef.TestEchoACK{}, func(ses cellnet.CellID, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("server recv:", msg.String())

		cellnet.Send(ses, &coredef.TestEchoACK{
			Content: proto.String("world"),
		})
	})

	ltvsocket.SpawnAcceptor("127.0.0.1:8001", dispatcher.PeerHandler(disp))
}

func client() {

	disp := dispatcher.NewPacketDispatcher()

	dispatcher.RegisterMessage(disp, coredef.TestEchoACK{}, func(ses cellnet.CellID, content interface{}) {
		msg := content.(*coredef.TestEchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	dispatcher.RegisterMessage(disp, coredef.ConnectedACK{}, func(ses cellnet.CellID, content interface{}) {
		cellnet.Send(ses, &coredef.TestEchoACK{
			Content: proto.String("hello"),
		})
	})

	ltvsocket.SpawnConnector("127.0.0.1:8001", dispatcher.PeerHandler(disp))

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
