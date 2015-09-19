package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type Session interface {
	Send(*cellnet.Packet)

	Close()
}

type ltvSession struct {
	writeChan chan *cellnet.Packet
	stream    cellnet.PacketStream

	OnClose func()
}

func (self *ltvSession) Send(pkt *cellnet.Packet) {

	if pkt == nil {
		return
	}

	self.writeChan <- pkt
}

func (self *ltvSession) Close() {

	// 关闭socket, 触发读错误
	self.stream.Close()
}

func newSession(stream cellnet.PacketStream, evq *cellnet.EvQueue) *ltvSession {

	self := &ltvSession{
		writeChan: make(chan *cellnet.Packet),
		stream:    stream,
	}

	var endSync sync.WaitGroup

	// 发送线程
	go func() {

		for {

			pkt := <-self.writeChan

			if pkt == nil {
				break
			}

			stream.Write(pkt)

		}

		// 告诉接收线程, 搞定
		endSync.Done()

	}()

	// 接收线程
	go func() {
		var err error
		var pkt *cellnet.Packet

		for {

			// 从Socket读取封包
			pkt, err = stream.Read()

			if err != nil {

				if self.OnClose != nil {
					self.OnClose()
				}

				evq.Post(NewDataEvent(Event_Closed, self, nil))
				break
			}

			evq.Post(&DataEvent{
				Packet: pkt,
				Ses:    self,
			})

		}

		// 通知发送线程停止
		self.writeChan <- nil

		// 置信号量
		endSync.Add(1)

		// 等待结束
		endSync.Wait()

	}()

	return self
}
