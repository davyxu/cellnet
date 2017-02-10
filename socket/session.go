package socket

import (
	"sync"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
)

type SocketSession struct {
	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭

	// handler相关上下文
	stream *ltvStream

	sendList *PacketList
}

func (self *SocketSession) ID() int64 {
	return self.id
}

func (self *SocketSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *SocketSession) Close() {
	self.sendList.Add(&cellnet.Packet{})
}

func (self *SocketSession) Send(data interface{}) {

	pkt, _ := cellnet.BuildPacket(data)

	msgLog("send", self, pkt)

	self.RawSend(pkt)
}

func (self *SocketSession) RawSend(pkt *cellnet.Packet) {

	if pkt != nil {
		self.sendList.Add(pkt)
	}
}

// 发送线程
func (self *SocketSession) sendThread() {

	var writeList []*cellnet.Packet

	for {

		willExit := false
		writeList = writeList[0:0]

		// 复制出队列
		packetList := self.sendList.BeginPick()

		for _, p := range packetList {

			if p.MsgID == 0 {
				willExit = true
			} else {
				writeList = append(writeList, p)
			}
		}

		self.sendList.EndPick()

		// 写队列
		for _, p := range writeList {

			if err := self.stream.Write(p); err != nil {
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

func (self *SocketSession) recvThread2(eq cellnet.EventQueue) {

	for {

		ev := NewSessionEvent(0, self, nil)

		if self.p.GetHandler().Call(SessionEvent_Recv, ev) != nil {
			break
		}
	}

	if self.needNotifyWrite {
		self.Close()
	}

	// 通知接收线程ok
	self.endSync.Done()
}

// 接收线程
func (self *SocketSession) recvThread(eq cellnet.EventQueue) {
	var err error
	var pkt *cellnet.Packet

	for {

		// 从Socket读取封包
		pkt, err = self.stream.Read()

		if err != nil {

			ev := newSessionEvent(Event_SessionClosed, self, &gamedef.SessionClosed{Reason: err.Error()})

			msgLog("recv", self, ev.Packet)

			// 断开事件
			eq.Post(self.p, ev)
			break
		}

		// 消息日志要多损耗一次解析性能

		msgLog("recv", self, pkt)

		// 逻辑封包
		eq.Post(self.p, &SessionEvent{
			Packet: pkt,
			Ses:    self,
		})

	}

	if self.needNotifyWrite {
		self.Close()
	}

	// 通知接收线程ok
	self.endSync.Done()
}

func newSession(stream *ltvStream, eq cellnet.EventQueue, p cellnet.Peer) *SocketSession {

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
	go self.recvThread2(eq)

	// 发送线程
	go self.sendThread()

	return self
}
