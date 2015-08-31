package nexus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/dispatcher"
	"github.com/davyxu/cellnet/ltvsocket"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/golang/protobuf/proto"
	"log"
)

func listenNexus() {

	listenAddr := config.Listen

	ltvsocket.SpawnAcceptor(listenAddr, dispatcher.PeerHandler(Dispatcher))

	dispatcher.RegisterMessage(Dispatcher, coredef.RegionLinkREQ{}, func(ses cellnet.CellID, content interface{}) {

		msg := content.(*coredef.RegionLinkREQ)

		profile := msg.GetProfile()

		if profile.GetID() == cellnet.RegionID {
			log.Printf("duplicate regionid: %d@%s", profile.GetID(), profile.GetAddress())
			return
		}

		addRegion(ses, profile)

		ack := coredef.RegionLinkACK{
			AddressList: make([]*coredef.Region, 0),
			Status: &coredef.Region{
				ID:      proto.Int32(cellnet.RegionID),
				Address: proto.String(config.Listen),
			},
		}

		IterateRegion(func(profile *RegionData) {

			ack.AddressList = append(ack.AddressList, profile.Region)

		})

		cellnet.Send(ses, cellnet.BuildPacket(&ack))

	})

	dispatcher.RegisterMessage(Dispatcher, coredef.ClosedACK{}, func(ses cellnet.CellID, _ interface{}) {

		removeRegion(ses)

	})
}
