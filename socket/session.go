package socket

import (
	"github.com/davyxu/cellnet"
	"net"
	"sync"
	"time"
)

type ltvSession struct {
	writeChan chan interface{}
	stream    PacketStream

	OnClose func()

	id uint32

	p cellnet.Peer

	endSync sync.WaitGroup
}

func (self *ltvSession) ID() uint32 {
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

// 接收线程
func (self *ltvSession) recvThread(evq *cellnet.EvQueue) {
	var err error
	var pkt *cellnet.Packet

	for {

		// 从Socket读取封包
		pkt, err = self.stream.Read()

		if err != nil {

			// 断开事件
			evq.Post(NewDataEvent(Event_Closed, self, nil))
			break
		}

		// 逻辑封包
		evq.Post(&DataEvent{
			Packet: pkt,
			Ses:    self,
		})

	}

	// 通知发送线程停止
	self.writeChan <- closeWrite{}

	// 通知接收线程ok
	self.endSync.Done()
}

// 发送线程
func (self *ltvSession) sendThread() {

	for {

		switch data := (<-self.writeChan).(type) {
		case closeWrite:
			goto exitsendloop
		default:
			if err := self.stream.Write(cellnet.BuildPacket(data)); err != nil {
				goto exitsendloop
			}
		}

	}

exitsendloop:

	// 通知发送线程ok
	self.endSync.Done()
}

// 退出线程
func (self *ltvSession) exitThread() {

	// 布置接收和发送2个任务
	self.endSync.Add(2)

	// 等待2个任务结束
	self.endSync.Wait()

	// 在这里断开session与逻辑的所有关系
	if self.OnClose != nil {
		self.OnClose()
	}

	// 延时断开
	time.AfterFunc(time.Second, func() {
		self.stream.Close()
	})
}

type closeWrite struct {
}

func newSession(stream PacketStream, evq *cellnet.EvQueue, p cellnet.Peer) *ltvSession {

	self := &ltvSession{
		writeChan: make(chan interface{}),
		stream:    stream,
		p:         p,
	}

	// 接收线程
	go self.recvThread(evq)

	// 发送线程
	go self.sendThread()

	// 退出线程
	go self.exitThread()

	return self
}
