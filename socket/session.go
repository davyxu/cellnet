package socket

import (
	"sync"

	"github.com/davyxu/cellnet"
	"github.com/golang/protobuf/proto"
)

type closeWritePacket struct {
}

type ltvSession struct {
	stream PacketStream

	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭

	sendList *PacketList
}

func (self *ltvSession) ID() int64 {
	return self.id
}

func (self *ltvSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *ltvSession) Close() {
	self.sendList.Add(&cellnet.Packet{})
}

func (self *ltvSession) Send(data interface{}) {

	pkt, meta := cellnet.BuildPacket(data)

	if EnableMessageLog {
		msgLog(&MessageLogInfo{
			Dir:       "send",
			PeerName:  self.FromPeer().Name(),
			SessionID: self.ID(),
			Name:      meta.Name,
			ID:        meta.ID,
			Size:      int32(len(pkt.Data)),
			Data:      data.(proto.Message).String(),
		})

	}

	self.RawSend(pkt)
}

func (self *ltvSession) RawSend(pkt *cellnet.Packet) {

	if pkt != nil {
		self.sendList.Add(pkt)
	}
}

// 发送线程
func (self *ltvSession) sendThread() {

	for {
		packetList := self.sendList.BeginPick()

		willExit := false

		for _, p := range packetList {

			if p.MsgID == 0 {
				willExit = true
			} else {

				if err := self.stream.Write(p); err != nil {
					willExit = true
					break
				}

			}

		}

		self.sendList.EndPick()

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

// 接收线程
func (self *ltvSession) recvThread(eq cellnet.EventQueue) {
	var err error
	var pkt *cellnet.Packet

	for {

		// 从Socket读取封包
		pkt, err = self.stream.Read()

		if err != nil {

			// 断开事件
			eq.PostData(NewSessionEvent(Event_SessionClosed, self, nil))
			break
		}

		// 逻辑封包
		eq.PostData(&SessionEvent{
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

func newSession(stream PacketStream, eq cellnet.EventQueue, p cellnet.Peer) *ltvSession {

	self := &ltvSession{
		stream:          stream,
		p:               p,
		needNotifyWrite: true,
		sendList:        NewPacketList(),
	}

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
