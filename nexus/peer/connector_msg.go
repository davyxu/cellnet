package peer

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/nexus/addrlist"
	"github.com/davyxu/cellnet/nexus/config"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

func joinNexus(addr string) {

	ltvsocket.SpawnConnector(addr, dispatcher.PeerHandler(Dispatcher))
	log.Printf("begin join: %s", addr)

	Dispatcher.RegisterCallback(dispatcher.MsgNewSession, func(src cellnet.CellID, _ *cellnet.Packet) {
		cellnet.Send(src, &coredef.RegionLinkREQ{
			Profile: &coredef.Region{
				ID:      proto.Int32(cellnet.RegionID),
				Address: proto.String(config.Data.ListenAddress),
			},
		})
	})

	dispatcher.RegisterMessage(Dispatcher, coredef.RegionLinkACK{}, func(src cellnet.CellID, content interface{}) {

		msg := content.(*coredef.RegionLinkACK)

		status := msg.GetStatus()

		if status.GetID() == cellnet.RegionID {
			log.Printf("duplicate regionid: %d@%s", status.GetID(), status.GetAddress())
			return
		}

		addrlist.AddRegion(src, status)

		for _, rg := range msg.GetAddressList() {

			log.Printf("address: %d@%s", rg.GetID(), rg.GetAddress())

			// 不能是自己
			if rg.GetID() == cellnet.RegionID {
				continue
			}

			// 已经连上了, 不再连接
			if addrlist.GetRegion(rg.GetID()) != nil {
				continue
			}

			// 连接地址中的服务器
			joinNexus(rg.GetAddress())
		}

	})

	//	Dispatcher.RegisterCallback(comm.EventClosed, func(src cellnet.CellID, _ *cellnet.Packet) {

	//		addrlist.RemoveRegion(src)

	//	})

}
