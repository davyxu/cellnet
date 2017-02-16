package socket

import (
	"sync"

	"github.com/davyxu/cellnet"
)

type SocketSession struct {
	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭

	// handler相关上下文
	stream *PacketStream

	sendList *eventList
}

func (self *SocketSession) ID() int64 {
	return self.id
}

func (self *SocketSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *SocketSession) Close() {
	self.sendList.Add(nil)
}

func (self *SocketSession) Send(data interface{}) {

	ev := cellnet.NewSessionEvent(cellnet.SessionEvent_Send, self)
	ev.Msg = data

	self.RawSend(ev.SendHandler, ev)

}

func (self *SocketSession) RawSend(sendHandler cellnet.EventHandler, ev *cellnet.SessionEvent) {

	if sendHandler == nil {
		_, sendHandler = self.p.GetHandler()
	}

	ev.Ses = self

	cellnet.HandlerCallFirst(sendHandler, ev)
}

// 发送线程
func (self *SocketSession) sendThread() {

	var writeList []*cellnet.SessionEvent

	for {

		willExit := false
		writeList = writeList[0:0]

		// 复制出队列
		packetList := self.sendList.BeginPick()

		for _, ev := range packetList {

			if ev == nil {
				willExit = true
				break
			} else {
				writeList = append(writeList, ev)
			}
		}

		self.sendList.EndPick()

		// 写队列
		for _, ev := range writeList {

			if err := self.stream.Write(ev); err != nil {
				willExit = true
				break
			}

		}

		if err := self.stream.Flush(); err != nil {
			willExit = true
		}

		if willExit {
			goto exitsendloop
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

func (self *SocketSession) recvThread(eq cellnet.EventQueue) {

	for {

		ev := cellnet.NewSessionEvent(cellnet.SessionEvent_Recv, self)

		recv, _ := self.p.GetHandler()

		if cellnet.HandlerCallFirst(recv, ev) != nil {
			break
		}
	}

	if self.needNotifyWrite {
		self.Close()
	}

	// 通知接收线程ok
	self.endSync.Done()
}

//// 接收线程
//func (self *SocketSession) recvThread(eq cellnet.EventQueue) {
//	var err error
//	var pkt *cellnet.Packet

//	for {

//		// 从Socket读取封包
//		pkt, err = self.stream.Read()

//		if err != nil {

//			ev := newSessionEvent(Event_SessionClosed, self, &gamedef.SessionClosed{Reason: err.Error()})

//			msgLog("recv", self, ev.Packet)

//			// 断开事件
//			eq.Post(self.p, ev)
//			break
//		}

//		// 消息日志要多损耗一次解析性能

//		msgLog("recv", self, pkt)

//		// 逻辑封包
//		eq.Post(self.p, &SessionEvent{
//			Packet: pkt,
//			Ses:    self,
//		})

//	}

//	if self.needNotifyWrite {
//		self.Close()
//	}

//	// 通知接收线程ok
//	self.endSync.Done()
//}

func newSession(stream *PacketStream, eq cellnet.EventQueue, p cellnet.Peer) *SocketSession {

	self := &SocketSession{
		stream:          stream,
		p:               p,
		needNotifyWrite: true,
		sendList:        NewPacketList(),
	}

	// 使用peer的统一设置
	self.stream.maxPacketSize = p.MaxPacketSize()

	// 布置接收和发送2个任务
	// bug fix感谢viwii提供的线索
	self.endSync.Add(2)

	go func() {

		// 等待2个任务结束
		self.endSync.Wait()

		// 在这里断开session与逻辑的所有关系
		if self.OnClose != nil {
			self.OnClose()
		}
	}()

	// 接收线程
	go self.recvThread(eq)

	// 发送线程
	go self.sendThread()

	return self
}
