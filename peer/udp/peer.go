package udp

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellpeer "github.com/davyxu/cellnet/peer"
	cellqueue "github.com/davyxu/cellnet/queue"
	xframe "github.com/davyxu/x/frame"
)

type Peer struct {
	cellpeer.Hooker
	*cellpeer.SessionManager
	xframe.PropertySet
	Queue *cellqueue.Queue
	Recv  func(ses *Session, data []byte) (ev *cellevent.RecvMsgEvent, err error)
	Send  func(ses *Session, ev *cellevent.SendMsgEvent) error
}

func (self *Peer) Peer() *Peer {
	return self
}

func newPeer() *Peer {
	return &Peer{
		SessionManager: cellpeer.NewSessionManager(),
	}
}

func PeerFromSession(ses cellnet.Session) *Peer {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		return tcpSes.peer
	}

	return nil
}

//func ConnectorFromSession(ses cellnet.Session) *Connector {
//	if ses == nil {
//		return nil
//	}
//
//	if tcpSes, ok := ses.(*Session); ok {
//		if conn, connOK := tcpSes.parent.(*Connector); connOK {
//			return conn
//		}
//	}
//
//	return nil
//}

func AcceptorFromSession(ses cellnet.Session) *Acceptor {
	if ses == nil {
		return nil
	}
	if tcpSes, ok := ses.(*Session); ok {
		if acc, accOK := tcpSes.parent.(*Acceptor); accOK {
			return acc
		}
	}

	return nil
}

func SessionQueuedCall(ses cellnet.Session, callback func()) {
	peer := PeerFromSession(ses)
	if peer == nil {
		return
	}

	if peer.Queue == nil {
		callback()
	} else {
		peer.Queue.Post(callback)
	}
}
