package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

func SpawnSession(stream cellnet.IPacketStream, callback func(cellnet.CellID, interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io goroutine
	go func() {
		var err error
		var pkt *cellnet.Packet

		for {

			// Read packet data as ltv packet format
			pkt, err = stream.Read()

			if err != nil {

				cellnet.Send(cid, err)
				break
			}

			cellnet.Send(cid, pkt)

		}

	}()

	return cid
}
