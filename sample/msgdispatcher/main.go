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

	dispatcher.RegisterMessage(disp, coredef.EchoACK{}, func(ses cellnet.CellID, rawmsg interface{}) {
		msg := rawmsg.(*coredef.EchoACK)

		log.Println("server recv:", msg.String())

		cellnet.Send(ses, &coredef.EchoACK{
			Content: proto.String("world"),
		})
	})

	ltvsocket.SpawnAcceptor("127.0.0.1:8001", dispatcher.PeerHandler(disp))
}

func client() {

	disp := dispatcher.NewPacketDispatcher()

	dispatcher.RegisterMessage(disp, coredef.EchoACK{}, func(ses cellnet.CellID, rawmsg interface{}) {
		msg := rawmsg.(*coredef.EchoACK)

		log.Println("client recv:", msg.String())

		done <- true
	})

	disp.RegisterCallback(dispatcher.EventNewSession, func(ses cellnet.CellID, _ *cellnet.Packet) {
		cellnet.Send(ses, &coredef.EchoACK{
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
