package nexus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

func joinNexus(addr string) {

	ltvsocket.SpawnConnector(addr, dispatcher.PeerHandler(Dispatcher))

	dispatcher.RegisterMessage(Dispatcher, coredef.ConnectedACK{}, func(src cellnet.CellID, _ interface{}) {
		cellnet.Send(src, cellnet.BuildPacket(&coredef.RegionLinkREQ{
			Profile: &coredef.Region{
				ID:      proto.Int32(cellnet.RegionID),
				Address: proto.String(config.Listen),
			},
		}))
	})

	dispatcher.RegisterMessage(Dispatcher, coredef.RegionLinkACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.RegionLinkACK)

		status := msg.GetStatus()

		if status.GetID() == cellnet.RegionID {
			log.Printf("[nexus] duplicate regionid: %d@%s", status.GetID(), status.GetAddress())
			return
		}

		addRegion(src, status)

		for _, rg := range msg.GetAddressList() {

			//log.Printf("address: %d@%s", rg.GetID(), rg.GetAddress())

			// 不能是自己
			if rg.GetID() == cellnet.RegionID {
				continue
			}

			// 已经连上了, 不再连接
			if GetRegion(rg.GetID()) != nil {
				continue
			}

			// 连接地址中的服务器
			joinNexus(rg.GetAddress())
		}

	})

}
