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
	ltvsocket.SpawnAcceptor("127.0.0.1:8001", func(mailbox chan interface{}) {
		for {

			switch v := (<-mailbox).(type) {
			case cellnet.IPacketStream:

				ltvsocket.SpawnSession(v, func(sesmail chan interface{}) {

					for {

						switch pkt := (<-sesmail).(type) {
						case *cellnet.Packet:

							log.Println("server recv:", cellnet.ReflectContent(pkt))

							v.Write(cellnet.BuildPacket(&coredef.EchoACK{
								Content: proto.String("world"),
							}))
						}

					}

				})

			case IError:
				log.Println(cellnet.ReflectContent(v))
			}
		}

	})
}

func client() {

	ltvsocket.SpawnConnector("127.0.0.1:8001", func(mailbox chan interface{}) {
		for {

			switch v := (<-mailbox).(type) {
			case cellnet.IPacketStream:

				// new session
				ltvsocket.SpawnSession(v, func(sesmail chan interface{}) {

					switch pkt := (<-sesmail).(type) {
					case *cellnet.Packet:

						var ack coredef.EchoACK
						if err := proto.Unmarshal(pkt.Data, &ack); err != nil {
							log.Println(err)
						} else {
							log.Println("client recv", ack.String())

							done <- true
						}

					}

				})

				// send packet on connected
				v.Write(cellnet.BuildPacket(&coredef.EchoACK{
					Content: proto.String("hello"),
				}))

			case IError:
				log.Println(cellnet.ReflectContent(v))

			}
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
