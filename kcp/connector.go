package kcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"github.com/xtaci/kcp-go"
)

// 连接器, 可由Peer转换
type KcpConnector interface {
	// 连接后的Session
	DefaultSession() cellnet.Session
	// 自动重连间隔, 0表示不重连, 默认不重连
	//SetAutoReconnectSec(sec int)
}

type kcpConnector struct {
	*kcpPeer
	defaultSes				cellnet.Session
}

func (k *kcpConnector) DefaultSession() cellnet.Session {
	return k.defaultSes
}

// 自动重连间隔=0不重连
//func (k *kcpConnector) SetAutoReconnectSec(sec int) {
//	k.autoReconnectSec = sec
//}

func (k *kcpConnector) Start(address string) cellnet.Peer {
	if k.IsRunning() {
		return k
	}
	k.SetAddress(address)
	k.SetRunning(true)
	defer k.SetRunning(false)
	//todo options
	kcpconn, err := kcp.DialWithOptions(k.Address(), nil, 10, 3)
	if err != nil {
		extend.PostSystemEvent(nil, cellnet.Event_ConnectFailed, k.ChainListRecv(), errToResult(err))
		return k
	}
	ses := newSession(kcpconn, k)
	k.defaultSes = ses
	k.Add(ses)
	ses.OnClose = func() {
		k.Remove(ses)
	}
	extend.PostSystemEvent(ses, cellnet.Event_Connected, k.ChainListRecv(), cellnet.Result_OK)
	ses.run()
	return  k
}

func (k *kcpConnector) Stop() {
	if !k.IsRunning() {
		return
	}
	if k.defaultSes != nil {
		k.defaultSes.Close()
	}
}

func NewConnector(q cellnet.EventQueue) cellnet.Peer {
	self := &kcpConnector{
		kcpPeer:				newKcpPeer(q, cellnet.NewSessionManager()),
	}
	return self
}

func errToResult(err error) cellnet.Result {
	if err == nil {
		return cellnet.Result_OK
	}
	return cellnet.Result_SocketError
}