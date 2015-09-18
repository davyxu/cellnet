package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

type ltvSession struct {
	writeChan chan *cellnet.Packet
}

func (self *ltvSession) Send(pkt *cellnet.Packet) {
	self.writeChan <- pkt
}

func newSession(stream cellnet.PacketStream, evq *cellnet.EvQueue) *ltvSession {

	self := &ltvSession{
		writeChan: make(chan *cellnet.Packet),
	}

	go func() {

		for {

			pkt := <-self.writeChan

			stream.Write(pkt)

		}

	}()

	// 接收线程
	go func() {
		var err error
		var pkt *cellnet.Packet

		for {

			// 从Socket读取封包
			pkt, err = stream.Read()

			if err != nil {

				evq.Post(NewDataEvent(Event_Closed, self, nil))
				break
			}

			evq.Post(&DataEvent{
				Packet: pkt,
				Ses:    self,
			})

		}

	}()

	return self
}
