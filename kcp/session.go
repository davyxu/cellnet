package kcp

import (
	"github.com/davyxu/cellnet"
	"io"
	"net"
)

type kcpSession struct {
	conn				net.Conn
	OnClose				func()
	id					int64
	p					cellnet.Peer
	tag					interface{}
	readChain 			*cellnet.HandlerChain
	writeChain			*cellnet.HandlerChain
}

func (k *kcpSession) RawConn() interface{} {
	return k.conn
}

func (k *kcpSession) Tag() interface{} {
	return k.tag
}

func (k *kcpSession) SetTag(tag interface{}) {
	k.tag = tag
}

func (k *kcpSession) ID() int64 {
	return k.id
}

func (k *kcpSession) SetID(id int64) {
	k.id = id
}

func (k *kcpSession) FromPeer() cellnet.Peer {
	return k.p
}

func (k *kcpSession) Close() {
	k.conn.Close()
}

func (k *kcpSession) DataSource() io.ReadWriter {
	return k.conn
}

func (k *kcpSession) Send(data interface{}) {
	ev := cellnet.NewEvent(cellnet.Event_Send, k)
	ev.Msg = data
	if ev.ChainSend == nil {
		ev.ChainSend = k.p.ChainSend()
	}
	k.RawSend(ev)
}

func (k *kcpSession) RawSend(ev *cellnet.Event) {
	if ev.Type != cellnet.Event_Send {
		panic("invalid event type, require Event_Send")
	}
	ev.Ses = k
	// 发送链处理: encode等操作
	if ev.ChainSend != nil {
		ev.ChainSend.Call(ev)
	}
	cellnet.MsgLog(ev)
	//// 写链处理
	k.writeChain.Call(ev)
}

func (k *kcpSession) recvThread() {
	for {
		ev := cellnet.NewEvent(cellnet.Event_Recv, k)
		k.readChain.Call(ev)
		// 接收日志
		cellnet.MsgLog(ev)
		k.p.ChainListRecv().Call(ev)
	}
}

func (k *kcpSession) run() {
	go k.recvThread()
}

func newSession(c net.Conn,p cellnet.Peer) *kcpSession {
	self := &kcpSession{
		conn:				c,
		p:					p,
	}
	self.readChain = p.CreateChainRead()
	self.writeChain = p.CreateChainWrite()
	return self
}
