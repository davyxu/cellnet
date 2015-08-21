package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"net"
)

func SpawnConnector(address string, callback func(cellnet.CellID, interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io goroutine
	go func() {

		conn, err := net.Dial("tcp", address)
		if err != nil {

			cellnet.Send(cid, err)
			return
		}

		cellnet.Send(cid, NewPacketStream(conn))

	}()

	return cid

}
