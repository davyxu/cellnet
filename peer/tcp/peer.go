package tcp

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
	cellpeer.SocketOption
	Queue *cellqueue.Queue
	Recv  func(ses *Session) (ev *cellevent.RecvMsgEvent, err error)
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

// SessionID根据各种实现不一样(例如网关), 应该在具体实现里获取
func SessionID(ses cellnet5.Session) int64 {
	if ses == nil {
		return 0
	}

	return ses.(interface {
		ID() int64
	}).ID()
}

func SessionPeer(ses cellnet5.Session) *Peer {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		return tcpSes.peer
	}

	return nil
}

func ConnectorFromSession(ses cellnet5.Session) *Connector {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		if conn, connOK := tcpSes.parent.(*Connector); connOK {
			return conn
		}
	}

	return nil
}

func AcceptorFromSession(ses cellnet5.Session) *Acceptor {
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

func SessionQueuedCall(ses cellnet5.Session, callback func()) {
	peer := SessionPeer(ses)
	if peer == nil {
		return
	}

	if peer.Queue == nil {
		callback()
	} else {
		peer.Queue.Post(callback)
	}
}
