package socket

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type closeWritePacket struct {
}

type RawSession interface {
	// 路由发包
	RelaySend(data interface{}, clientid int64)

	// 直接发送封包
	RawSend(*cellnet.Packet)

	cellnet.Session
}

type ltvSession struct {
	writeChan chan interface{}
	stream    PacketStream

	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭
}

func (self *ltvSession) ID() int64 {
	return self.id
}

func (self *ltvSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *ltvSession) Close() {

	self.writeChan <- closeWritePacket{}
}

func (self *ltvSession) Send(data interface{}) {

	self.RawSend(cellnet.BuildPacket(data))
}

func (self *ltvSession) RelaySend(data interface{}, clientid int64) {

	// TODO 如果性能不理想, 丢到线程中利用线程发送
	pkt := cellnet.BuildPacket(data)
	pkt.ClientID = clientid

	self.RawSend(pkt)

}

func (self *ltvSession) RawSend(pkt *cellnet.Packet) {

	if pkt == nil {
		return
	}

	self.writeChan <- pkt
}

// 发送线程
func (self *ltvSession) sendThread() {

	for {

		switch pkt := (<-self.writeChan).(type) {
		// 关闭循环
		case closeWritePacket:
			goto exitsendloop
		// 封包
		case *cellnet.Packet:
			if err := self.stream.Write(pkt); err != nil {
				goto exitsendloop
			}
		}

	}

exitsendloop:

	// 不需要读线程再次通知写线程
	self.needNotifyWrite = false

	// 关闭socket,触发读错误, 结束读循环
	self.stream.Close()

	// 通知发送线程ok
	self.endSync.Done()
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
			evq.Post(NewSessionEvent(Event_SessionClosed, self, nil))
			break
		}

		// 逻辑封包
		evq.Post(&SessionEvent{
			Packet: pkt,
			Ses:    self,
		})

	}

	if self.needNotifyWrite {

		// 通知发送线程停止
		self.writeChan <- closeWritePacket{}
	}

	// 通知接收线程ok
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

}

func newSession(stream PacketStream, evq *cellnet.EvQueue, p cellnet.Peer) *ltvSession {

	self := &ltvSession{
		writeChan:       make(chan interface{}),
		stream:          stream,
		p:               p,
		needNotifyWrite: true,
	}

	// 接收线程
	go self.recvThread(evq)

	// 发送线程
	go self.sendThread()

	// 退出线程
	go self.exitThread()

	return self
}
