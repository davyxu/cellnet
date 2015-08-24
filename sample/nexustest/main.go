package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/nexus/addrlist"
	"github.com/davyxu/cellnet/nexus/config"
	"github.com/davyxu/cellnet/nexus/express"
	"github.com/davyxu/cellnet/nexus/peer"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

var done = make(chan bool)

func server() {
	//  开始测试的时机
	addrlist.Event.Add("addrlist.AddRegion", func(args ...interface{}) {
		rd := args[0].(*addrlist.RegionData)

		// client 连上来了
		if rd.GetID() == 1 {
			log.Println("send")
			cellnet.Send(cellnet.NewCellID(1, 1), &coredef.TestEchoACK{
				Content: proto.String("send to node"),
			})

			cellnet.Send(cellnet.NewCellID(1, 4), &coredef.TestEchoACK{
				Content: proto.String("send to callback"),
			})
		}
	})
}

func client() {
	express.Register(peer.Dispatcher)

	dispatcher.RegisterMessage(peer.Dispatcher, coredef.TestEchoACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.TestEchoACK)
		log.Println("recv msg callback :", msg.GetContent())
	})

	cid := cellnet.Spawn(func(src cellnet.CellID, data interface{}) {
		log.Println("recv node msg", src.String(), data)
	})

	log.Println(cid.String())
}

func main() {

	if config.Data.TestCase == "express" {

		// 保证host先启动
		// 主机参数: -listen=127.0.0.1:7001 -test=express
		// 从机参数: -region=1 -test=express -listen=127.0.0.1:7002 -join=127.0.0.1:7001

		// host
		if cellnet.RegionID == 0 {

			server()

		} else { // client
			client()
		}

	}

	<-done

}
