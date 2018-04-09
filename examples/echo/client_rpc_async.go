package main

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/rpc"
	"time"
)

func clientAsyncRPC() {
	// 等待服务器返回数据
	done := make(chan struct{})

	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "async rpc", peerAddress, queue)

	// 创建一个消息同步接收器
	rv := proc.NewSyncReceiver(p)

	proc.BindProcessorHandler(p, "tcp.ltv", rv.EventCallback())

	p.Start()

	queue.StartLoop()

	// 等连接上时
	rv.WaitMessage("cellnet.SessionConnected")

	// 异步RPC
	rpc.Call(p, &TestEchoACK{
		Msg:   "hello",
		Value: 1234,
	}, time.Second, func(raw interface{}) {

		switch result := raw.(type) {
		case error:
			fmt.Println(result)
		default:
			fmt.Println(result)
			done <- struct{}{}
		}

	})

	// 等待客户端收到消息
	<-done
}
