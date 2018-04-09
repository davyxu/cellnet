package main

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/rpc"
)

func server() {
	queue := cellnet.NewEventQueue()

	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", peerAddress, queue)

	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted: // 接受一个连接
			fmt.Println("server accepted")
		case *TestEchoACK: // 收到连接发送的消息

			fmt.Printf("server recv %+v\n", msg)

			ack := &TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			}

			// 当服务器收到的是一个rpc消息
			if rpcevent, ok := ev.(*rpc.RecvMsgEvent); ok {

				// 以RPC方式回应
				rpcevent.Reply(ack)
			} else {

				// 收到的是普通消息，回普通消息
				ev.Session().Send(ack)
			}

		case *cellnet.SessionClosed: // 连接断开
			fmt.Println("session closed: ", ev.Session().ID())
		}

	})

	peerIns.Start()

	queue.StartLoop()
}
