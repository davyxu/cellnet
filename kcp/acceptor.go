package kcp

import (
	kcp "github.com/xtaci/kcp-go"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/extend"
	"net"
)

type kcpAcceptor struct {
	*kcpPeer
	listener				net.Listener

}

func (k *kcpAcceptor) Start(address string) cellnet.Peer {
	if k.IsRunning() {
		return k
	}
	k.SetAddress(address)
	var err error
	k.listener, err = kcp.ListenWithOptions(address, nil, 10, 3)
	if err != nil {
		log.Errorf("#listen failed(%s) %v", k.NameOrAddress(), err.Error())
		return k
	}
	log.Infof("#listen(%s) %s", k.Name(), k.Address())
	go k.accept()
	return k
}

func (k *kcpAcceptor) accept() {
	k.SetRunning(true)
	for {
		conn, err := k.listener.Accept()
		if err != nil {
			if log.IsDebugEnabled() {
				log.Errorf("#accept failed(%s) %v", k.NameOrAddress(), err.Error())
			}
			extend.PostSystemEvent(nil, cellnet.Event_AcceptFailed, k.ChainListRecv(), errToResult(err))
			break
		}
		ses := newSession(conn, k)
		k.Add(ses)
		ses.OnClose = func() {
			k.Remove(ses)
		}
		// 投递连接已接受事件
		extend.PostSystemEvent(ses, cellnet.Event_Accepted, k.ChainListRecv(), cellnet.Result_OK)
		// 事件处理完成开始处理数据收发
		ses.run()
	}
	k.SetRunning(false)
}

func (k *kcpAcceptor) Stop() {
	if !k.IsRunning() {
		return
	}
	k.listener.Close()
	k.CloseAllSession()
}

func NewAcceptor(q cellnet.EventQueue) cellnet.Peer {
	self := &kcpAcceptor{
		kcpPeer:				newKcpPeer(q, cellnet.NewSessionManager()),
	}
	return self
}
