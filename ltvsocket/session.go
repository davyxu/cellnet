package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

func SpawnSession(stream cellnet.IPacketStream, callback func(cellnet.CellID, interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io线程
	go func() {
		var err error
		var pkt *cellnet.Packet

		for {

			// 从Socket读取封包并转为ltv格式
			pkt, err = stream.Read()

			if err != nil {

				cellnet.Send(cid, EventClose{error: err})
				break
			}

			cellnet.Send(cid, pkt)

		}

	}()

	return cid
}
