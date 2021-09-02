package tcp

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellpeer "github.com/davyxu/cellnet/peer"
	cellqueue "github.com/davyxu/cellnet/queue"
	"github.com/davyxu/x/container"
)

type Peer struct {
	cellpeer.Hooker
	*cellpeer.SessionManager
	xcontainer.Mapper
	cellpeer.SocketOption
	cellpeer.Protect
	Queue  *cellqueue.Queue
	OnRecv func(ses *Session) (ev *cellevent.RecvMsg, err error)
	OnSend func(ses *Session, ev *cellevent.SendMsg) error
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
func SessionID(ses cellnet.Session) int64 {
	if ses == nil {
		return 0
	}

	return ses.(interface {
		ID() int64
	}).ID()
}

func SessionRaw(ses cellnet.Session) *Session {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		return tcpSes
	}

	return nil
}

func SessionPeer(ses cellnet.Session) *Peer {
	raw := SessionRaw(ses)
	if raw == nil {
		return nil
	}

	return raw.Peer
}

func SessionQueuedCall(ses cellnet.Session, callback func()) {
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
