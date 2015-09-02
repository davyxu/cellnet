package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

func SpawnAcceptor(address string, callback func(interface{})) cellnet.CellID {

	cid := cellnet.Spawn(callback)

	// io goroutine
	go func() {

		if config.SocketLog {
			log.Printf("[socket] #listen %s %s\n", cid.String(), address)
		}

		ln, err := net.Listen("tcp", address)

		if err != nil {
			cellnet.Send(cid, SocketListenError{Err: err})

			if config.SocketLog {
				log.Println("[socket] listen failed", err.Error())
			}

			return
		}

		for {
			conn, err := ln.Accept()

			if err != nil {
				continue
			}

			cellnet.Send(cid, SocketCreateSession{Stream: NewPacketStream(conn), Type: SessionAccepted})
		}

	}()

	return cid

}
