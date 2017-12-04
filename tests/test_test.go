package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/packet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/tests/proto"
	"reflect"
	"testing"
)

const testAddress = "127.0.0.1:8001"

func server() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		TypeName: "tcp.Acceptor",
		Queue:    queue,
		Address:  testAddress,
		Name:     "server",
		Event: packet.NewMessageCallback(func(ses cellnet.Session, raw interface{}) {
			switch ev := raw.(type) {
			case socket.AcceptedEvent:
				fmt.Println("server accepted")
			case packet.MsgEvent:

				msg := ev.Msg.(*proto.TestEchoACK)

				fmt.Printf("server recv %+v\n", msg)

				ses.Send(&proto.TestEchoACK{
					Msg:   "ack",
					Value: msg.Value,
				})
			case socket.SessionClosedEvent:
				fmt.Println("server error: ", ev.Error)
			}
		}),
	}).Start()

	queue.StartLoop()

	queue.Wait()
}

func client() {
	queue := cellnet.NewEventQueue()

	cellnet.NewPeer(cellnet.PeerConfig{
		TypeName: "tcp.Connector",
		Queue:    queue,
		Address:  testAddress,
		Name:     "client",
		Event: packet.NewMessageCallback(func(ses cellnet.Session, raw interface{}) {

			switch ev := raw.(type) {
			case socket.ConnectedEvent:
				fmt.Println("client connected")
				ses.Send(&proto.TestEchoACK{
					Msg:   "hello",
					Value: 1234,
				})
			case packet.MsgEvent:

				msg := ev.Msg.(*proto.TestEchoACK)

				fmt.Printf("client recv %+v\n", msg)

				queue.StopLoop(0)
			case socket.SessionClosedEvent:
				fmt.Println("client error: ", ev.Error)
			}

		}),
	}).Start()

	queue.StartLoop()

	queue.Wait()
}

func TestOneWayEcho(t *testing.T) {

	go server()

	client()
}

func TestMessageMeta(t *testing.T) {

	// 打印消息名
	fmt.Println(cellnet.MessageMetaByID(1).Name)

	// 通过消息反射类型，查询ID
	fmt.Println(cellnet.MessageMetaByType(reflect.TypeOf(&proto.TestEchoACK{}).Elem()).ID)

	// 使用消息名，获取消息类型
	fmt.Println(cellnet.MessageMetaByName("test.TestEchoACK").Type.Name())
}
