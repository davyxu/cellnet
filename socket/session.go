package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"sync"
	"time"
)

type socketSession struct {
	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭

	// handler相关上下文
	stream cellnet.PacketStream

	sendList *eventList
}

func (self *socketSession) Stream() cellnet.PacketStream {
	return self.stream
}

func (self *socketSession) ID() int64 {
	return self.id
}

func (self *socketSession) SetID(id int64) {
	self.id = id
}

func (self *socketSession) FromPeer() cellnet.Peer {
	return self.p
}

func (self *socketSession) Close() {
	self.sendList.Add(nil)
}

func (self *socketSession) Send(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Send, self)
	ev.Msg = data

	self.RawSend(ev.SendHandler, ev)

}

func (self *socketSession) RawSend(sendHandler []cellnet.EventHandler, ev *cellnet.Event) {

	if sendHandler == nil {
		_, sendHandler = self.p.HandlerList()
	}

	ev.Ses = self

	cellnet.HandlerChainCall(sendHandler, ev)
}

func (self *socketSession) Post(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Post, self)

	ev.Msg = data

	cellnet.MsgLog(ev)

	self.p.Call(ev)
}

func (self *socketSession) RawPost(recvHandler []cellnet.EventHandler, ev *cellnet.Event) {
	if recvHandler == nil {
		recvHandler, _ = self.p.HandlerList()
	}

	ev.Ses = self

	cellnet.HandlerChainCall(recvHandler, ev)
}

// 发送线程
func (self *socketSession) sendThread() {

	for {

		// 写超时
		_, write := self.FromPeer().(SocketOptions).SocketDeadline()

		if write != 0 {
			self.stream.Raw().SetWriteDeadline(time.Now().Add(write))
		}

		writeList, willExit := self.sendList.Pick()

		// 写队列
		for _, ev := range writeList {

			if err := self.stream.Write(ev.MsgID, ev.Data); err != nil {
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

func (self *socketSession) recvThread() {

	// 暂时不支持运行期修改HandlerList
	recvList, _ := self.p.HandlerList()

	for {

		ev := cellnet.NewEvent(cellnet.Event_Recv, self)

		cellnet.HandlerChainCall(recvList, ev)

		if ev.Result() != cellnet.Result_OK {
			extend.PostSystemEvent(ev.Ses, cellnet.Event_Closed, recvList, ev.Result())
			break
		}

	}

	if self.needNotifyWrite {
		self.Close()
	}

	// 通知接收线程ok
	self.endSync.Done()
}

func (self *socketSession) run() {
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
	go self.recvThread()

	// 发送线程
	go self.sendThread()
}

func newSession(stream cellnet.PacketStream, p cellnet.Peer) *socketSession {

	self := &socketSession{
		stream:          stream,
		p:               p,
		needNotifyWrite: true,
		sendList:        NewPacketList(),
	}

	// 使用peer的统一设置
	if s, ok := self.stream.(*TLVStream); ok {

		s.SetMaxPacketSize(p.(SocketOptions).MaxPacketSize())
	}

	return self
}
