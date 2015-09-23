package socket

import (
	"github.com/davyxu/cellnet"
	"net"
	"sync"
	"time"
)

type ltvSession struct {
	writeChan chan interface{}
	stream    cellnet.PacketStream

	OnClose func()

	id int64

	p cellnet.Peer
}

func (self *ltvSession) SetID(id int64) {
	self.id = id
}

func (self *ltvSession) ID() int64 {
	return self.id
}

func (self *ltvSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *ltvSession) Send(data interface{}) {

	if data == nil {
		return
	}

	self.writeChan <- data
}

func (self *ltvSession) Close() {

	tcpConn := self.stream.Raw().(*net.TCPConn)

	// 关闭socket, 触发读错误
	tcpConn.CloseRead()
}

type closeWrite struct {
}

func newSession(stream cellnet.PacketStream, evq *cellnet.EvQueue, p cellnet.Peer) *ltvSession {

	self := &ltvSession{
		writeChan: make(chan interface{}),
		stream:    stream,
		p:         p,
	}

	var endSync sync.WaitGroup
	endSync.Add(2)

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

		// 通知发送线程停止
		self.writeChan <- closeWrite{}

		// 等待结束
		endSync.Done()

	}()

	// 发送线程
	go func() {

		for {

			switch data := (<-self.writeChan).(type) {
			case closeWrite:
				goto ExitWriteLoop
			default:
				if err := stream.Write(cellnet.BuildPacket(data)); err != nil {
					goto ExitWriteLoop
				}
			}

		}

	ExitWriteLoop:

		// 告诉接收线程, 搞定
		endSync.Done()

	}()

	go func() {

		endSync.Wait()

		if self.OnClose != nil {
			self.OnClose()
		}

		// 延时断开
		time.AfterFunc(time.Second, func() {
			stream.Close()
		})

	}()

	return self
}
