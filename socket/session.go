package socket

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"io"
	"net"
	"sync"
	"time"
)

type socketSession struct {
	OnClose func() // 关闭函数回调

	id int64

	p cellnet.Peer

	endSync sync.WaitGroup

	needNotifyWrite bool // 是否需要通知写线程关闭

	sendList *eventList

	conn net.Conn

	tag interface{}

	readChain *cellnet.HandlerChain

	writeChain *cellnet.HandlerChain
}

func (self *socketSession) Tag() interface{} {
	return self.tag
}
func (self *socketSession) SetTag(tag interface{}) {
	self.tag = tag
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

func (self *socketSession) DataSource() io.ReadWriter {

	return self.conn
}

func (self *socketSession) Close() {
	self.sendList.Add(nil)
}

func (self *socketSession) Send(data interface{}) {

	ev := cellnet.NewEvent(cellnet.Event_Send, self)
	ev.Msg = data

	if ev.ChainSend == nil {
		ev.ChainSend = self.p.ChainSend()
	}

	self.RawSend(ev)

}

func (self *socketSession) RawSend(ev *cellnet.Event) {

	if ev.Type != cellnet.Event_Send {
		panic("invalid event type, require Event_Send")
	}

	ev.Ses = self

	if ev.ChainSend != nil {
		ev.ChainSend.Call(ev)
	}

	// 发送日志
	cellnet.MsgLog(ev)

	self.sendList.Add(ev)
}

func (self *socketSession) recvThread() {

	for {

		ev := cellnet.NewEvent(cellnet.Event_Recv, self)

		read, _ := self.FromPeer().(SocketOptions).SocketDeadline()

		if read != 0 {
			self.conn.SetReadDeadline(time.Now().Add(read))
		}

		self.readChain.Call(ev)

		if ev.Result() != cellnet.Result_OK {
			goto onClose
		}

		// 接收日志
		cellnet.MsgLog(ev)

		self.p.ChainListRecv().Call(ev)

		if ev.Result() != cellnet.Result_OK {
			goto onClose
		}

		continue

	onClose:
		extend.PostSystemEvent(ev.Ses, cellnet.Event_Closed, self.p.ChainListRecv(), ev.Result())
		break
	}

	if self.needNotifyWrite {
		self.Close()
	}

	// 通知接收线程ok
	self.endSync.Done()

}

// 发送线程
func (self *socketSession) sendThread() {

	for {

		// 写超时
		_, write := self.FromPeer().(SocketOptions).SocketDeadline()

		if write != 0 {
			self.conn.SetWriteDeadline(time.Now().Add(write))
		}

		writeList, willExit := self.sendList.Pick()

		// 写队列
		for _, ev := range writeList {

			self.writeChain.Call(ev)

			if ev.Result() != cellnet.Result_OK {
				willExit = true
			}

		}

		//if err := self.conn.Flush(); err != nil {
		//	willExit = true
		//}

		if willExit {
			goto exitsendloop
		}
	}

exitsendloop:

	// 不需要读线程再次通知写线程
	self.needNotifyWrite = false

	// 关闭socket,触发读错误, 结束读循环
	self.conn.Close()

	// 通知发送线程ok
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

func newSession(conn net.Conn, p cellnet.Peer) *socketSession {

	p.(interface {
		Apply(conn net.Conn)
	}).Apply(conn)

	self := &socketSession{
		conn:            conn,
		p:               p,
		needNotifyWrite: true,
		sendList:        NewPacketList(),
	}

	self.readChain = p.CreateChainRead()

	self.writeChain = p.CreateChainWrite()

	return self
}
