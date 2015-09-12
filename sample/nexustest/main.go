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

	//  开始测试的时机
	nexus.Event.Add("OnAddRegion", func(args ...interface{}) {
		rd := args[0].(*nexus.RegionData)

		// client 连上来了
		if rd.GetID() == 1 {
			log.Println("client connected", rd.GetID())
			cellnet.Send(cellnet.NewCellID(1, 3), cellnet.BuildPacket(&coredef.TestEchoACK{
				Content: proto.String("send to node"),
			}))

			//			cellnet.Send(cellnet.NewCellID(1, 4), &coredef.TestEchoACK{
			//				Content: proto.String("send to callback"),
			//			})
		}
	})
}

func client() {

	disp := dispatcher.NewDataDispatcher()

	// cellid: 1.3
	cid := cellnet.Spawn(func(data interface{}) {

		switch ev := data.(type) {

		case cellnet.SessionPacket: // 收

			disp.Call(int(ev.GetPacket().MsgID), data)
		default:
			log.Println("unknown data", cellnet.ReflectContent(ev))
		}

	})

	log.Println(cid.String())

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
