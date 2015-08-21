package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"net"
)

func SpawnAcceptor(address string, callback func(cellnet.CellID, interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io goroutine
	go func() {

		ln, err := net.Listen("tcp", address)

		if err != nil {
			cellnet.Send(cid, err.Error())
			return
		}

		for {
			conn, err := ln.Accept()

			if err != nil {
				continue
			}

			cellnet.Send(cid, NewPacketStream(conn))
		}

	}()

	return cid

}
