package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

type IError interface {
	Error() string
}

var done = make(chan bool)

func server() {
	ltvsocket.SpawnAcceptor("127.0.0.1:8001", func(self cellnet.CellID, cm interface{}) {

		switch v := cm.(type) {
		case ltvsocket.EventAccepted:

			ltvsocket.SpawnSession(v.Stream(), func(ses cellnet.CellID, sescm interface{}) {

				switch pkt := sescm.(type) {
				case *cellnet.Packet:

					log.Println("server recv:", cellnet.ReflectContent(pkt))

					v.Stream().Write(cellnet.BuildPacket(&coredef.TestEchoACK{
						Content: proto.String("world"),
					}))
				}

			})

		case IError:
			log.Println(cellnet.ReflectContent(v))
		}

	})
}

func client() {

	ltvsocket.SpawnConnector("127.0.0.1:8001", func(self cellnet.CellID, cm interface{}) {

		switch v := cm.(type) {
		case ltvsocket.EventConnected:

			// new session
			ltvsocket.SpawnSession(v.Stream(), func(ses cellnet.CellID, sescm interface{}) {

				switch pkt := sescm.(type) {
				case *cellnet.Packet:

					var ack coredef.TestEchoACK
					if err := proto.Unmarshal(pkt.Data, &ack); err != nil {
						log.Println(err)
					} else {
						log.Println("client recv", ack.String())

						done <- true
					}

				}

			})

			// send packet on connected
			v.Stream().Write(cellnet.BuildPacket(&coredef.TestEchoACK{
				Content: proto.String("hello"),
			}))

		case IError:
			log.Println(cellnet.ReflectContent(v))

		}

	})

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
