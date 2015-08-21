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

func listenNexus() {

	listenAddr := config.Data.ListenAddress

	ltvsocket.SpawnAcceptor(listenAddr, dispatcher.PeerHandler(Dispatcher))
	log.Printf("listen: %s", listenAddr)

	dispatcher.RegisterMessage(Dispatcher, coredef.RegionLinkREQ{}, func(ses cellnet.CellID, content interface{}) {

		msg := content.(*coredef.RegionLinkREQ)

		profile := msg.GetProfile()

		if profile.GetID() == cellnet.RegionID {
			log.Printf("duplicate regionid: %d@%s", profile.GetID(), profile.GetAddress())
			return
		}

		addrlist.AddRegion(ses, profile)

		ack := coredef.RegionLinkACK{
			AddressList: make([]*coredef.Region, 0),
			Status: &coredef.Region{
				ID:      proto.Int32(cellnet.RegionID),
				Address: proto.String(config.Data.ListenAddress),
			},
		}

		addrlist.IterateRegion(func(profile *addrlist.RegionData) {

			ack.AddressList = append(ack.AddressList, profile.Region)

		})

		cellnet.Send(ses, &ack)

	})

	Dispatcher.RegisterCallback(dispatcher.MsgClose, func(ses cellnet.CellID, _ interface{}) {

		addrlist.RemoveRegion(ses)

	})
}
