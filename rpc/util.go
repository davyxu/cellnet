package rpc

import (
	"errors"

	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

var (
	ErrInvalidPeerSession    = errors.New("rpc: Invalid peer type, require cellnet.RPCSessionGetter or cellnet.Session")
	ErrReplayMessageNotFound = errors.New("rpc: Reply message name not found")
)

type RPCSessionGetter interface {
	RPCSession() cellnet.Session
}

// 从peer获取rpc使用的session
func getPeerSession(ud interface{}) (cellnet.Session, error) {

	if ud == nil {
		return nil, ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		return i.RPCSession(), nil
	case cellnet.Session:
		return i, nil
	default:
		return nil, ErrInvalidPeerSession
	}
}

var (
	rpcIDSeq        int64
	requestByCallID sync.Map
)

type request struct {
	id     int64
	onRecv func(interface{})
}

var ErrTimeout = errors.New("time out")

func (self *request) RecvFeedback(msg interface{}) {
	self.onRecv(msg)
}

func (self *request) Send(ses cellnet.Session, msg interface{}) {
	data, msgid, _ := cellnet.EncodeMessage(msg)

	ses.Send(&RemoteCallREQ{
		MsgID:  msgid,
		Data:   data,
		CallID: self.id,
	})
}

func createRequest(onRecv func(interface{})) *request {

	self := &request{
		onRecv: onRecv,
	}

	self.id = atomic.AddInt64(&rpcIDSeq, 1)

	requestByCallID.Store(self.id, self)

	return self
}

func requestExists(callid int64) bool {

	_, ok := requestByCallID.Load(callid)
	return ok
}

func getRequest(callid int64) *request {

	if v, ok := requestByCallID.Load(callid); ok {

		requestByCallID.Delete(callid)
		return v.(*request)
	}

	return nil
}
