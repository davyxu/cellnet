package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/nexus"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

var done = make(chan bool)

func host() {
	// cellid: 0.1 acceptor
	// cellid: 0.2 connector

	//  开始测试的时机
	nexus.Event.Add("OnAddRegion", func(args ...interface{}) {
		rd := args[0].(*nexus.RegionData)

		// client 连上来了
		if rd.GetID() == 1 {
			log.Println("client connected", rd.GetID())
			cellnet.Send(cellnet.NewCellID(1, 3), &coredef.TestEchoACK{
				Content: proto.String("send to node"),
			})

			//			cellnet.Send(cellnet.NewCellID(1, 4), &coredef.TestEchoACK{
			//				Content: proto.String("send to callback"),
			//			})
		}
	})
}

func client() {
	// cellid: 1.1 acceptor
	// cellid: 1.2 connector

	// cellid: 1.3
	cid := cellnet.Spawn(func(src cellnet.CellID, data interface{}) {

		switch d := data.(type) {
		case *cellnet.Packet:
			log.Println("recv node msg", src.String(), cellnet.ReflectContent(d))
		}

	})

	log.Println(cid.String())

	nexus.Register(nexus.Dispatcher)

	// cellid: 1.4
	dispatcher.RegisterMessage(nexus.Dispatcher, coredef.TestEchoACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.TestEchoACK)
		log.Println("recv msg callback :", msg.GetContent())
	})

}

func main() {

	// 保证host先启动
	// 主机参数: ./host.toml
	// 从机参数: ./client.toml

	// host
	if cellnet.RegionID == 0 {

		host()

	} else { // client
		client()
	}

	<-done

}
