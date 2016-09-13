package socket

import (
	"sync"
	"time"

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

	sendList      []interface{}
	sendListGuard sync.Mutex
	sendListCond  *sync.Cond
}

func (self *ltvSession) ID() int64 {
	return self.id
}

func (self *ltvSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *ltvSession) Close() {
	self.pushSendMsg(closeWritePacket{})
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
		self.pushSendMsg(pkt)
	}
}

func (self *ltvSession) pushSendMsg(msg interface{}) {
	self.sendListGuard.Lock()
	self.sendList = append(self.sendList, msg)
	self.sendListGuard.Unlock()

	self.sendListCond.Signal()
}

// 发送线程
func (self *ltvSession) sendThread() {

	packetList := make([]*cellnet.Packet, 0)

	for {
		self.sendListGuard.Lock()

		for len(self.sendList) == 0 {
			self.sendListCond.Wait()
		}

		self.sendListGuard.Unlock()

		//强制sleep，积累消息(用于批量flush)
		time.Sleep(time.Millisecond)

		willExit := false

		self.sendListGuard.Lock()

		for _, v := range self.sendList {
			switch v.(type) {
			case *cellnet.Packet:
				packetList = append(packetList, v.(*cellnet.Packet))
			case closeWritePacket:
				willExit = true
			}
		}
		self.sendList = self.sendList[0:0]

		self.sendListGuard.Unlock()

		for i, p := range packetList {
			//当发送最后一个消息,且后续消息列表为空时，进行flush
			if err := self.stream.Write(p, len(packetList) == (i+1) && len(self.sendList) == 0); err != nil {
				willExit = true
				break
			}
		}

		packetList = packetList[0:0]

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
	}

	self.sendListCond = sync.NewCond(&self.sendListGuard)

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
