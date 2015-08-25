package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

func SpawnConnector(address string, callback func(cellnet.CellID, interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io goroutine
	go func() {

		if config.SocketLog {
			log.Printf("[socket] #connect %s %s\n", cid.String(), address)
		}

		conn, err := net.Dial("tcp", address)
		if err != nil {

			cellnet.Send(cid, EventConnectError{error: err})

			if config.SocketLog {
				log.Println("[socket] connect failed", err.Error())
			}
			return
		}

		cellnet.Send(cid, EventConnected{stream: NewPacketStream(conn)})

	}()

	return cid

}
