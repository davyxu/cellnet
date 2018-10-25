package rpc

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"sync"
	"sync/atomic"
)

var (
	rpcIDSeq        int64
	requestByCallID sync.Map
)

type request struct {
	id     int64
	onRecv func(interface{})
}

var ErrTimeout = errors.New("RPC time out")

func (self *request) RecvFeedback(msg interface{}) {

	// 异步和同步执行复杂，队列处理在具体的逻辑中手动处理
	self.onRecv(msg)
}

func (self *request) Send(ses cellnet.Session, msg interface{}) {

	//ctx, _ := ses.(cellnet.ContextSet)

	data, meta, err := codec.EncodeMessage(msg, nil)

	if err != nil {
		log.Errorf("rpc request message encode error: %s", err)
		return
	}

	ses.Send(&RemoteCallREQ{
		MsgID:  uint32(meta.ID),
		Data:   data,
		CallID: self.id,
	})

	//codec.FreeCodecResource(meta.Codec, data, ctx)
}

func createRequest(onRecv func(interface{})) *request {

	self := &request{
		onRecv: onRecv,
	}

	self.id = atomic.AddInt64(&rpcIDSeq, 1)

	requestByCallID.Store(self.id, self)

	return self
}

func getRequest(callid int64) *request {

	if v, ok := requestByCallID.Load(callid); ok {

		requestByCallID.Delete(callid)
		return v.(*request)
	}

	return nil
}
