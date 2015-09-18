package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"log"
	"net"
)

type ltvConnector struct {
	queue *cellnet.EvQueue
	conn  net.Conn
}

func (self *ltvConnector) Start(address string) {

	go func() {
		cn, err := net.Dial("tcp", address)

		if err != nil {

			log.Println("[socket] cononect failed", err.Error())
			return
		}

		self.conn = cn

		ses := newSession(NewPacketStream(cn), self.queue)

		self.queue.Post(NewDataEvent(Event_Connected, ses, nil))

	}()

}

func (self *ltvConnector) Stop() {

	if self.conn != nil {
		self.conn.Close()
	}

}

func newConnector(evq *cellnet.EvQueue) *ltvConnector {
	return &ltvConnector{queue: evq}
}
