package cellnet

import (
	"errors"
	"sync"
	"time"
)

type rpcResponse struct {
	pkt *Packet

	callid int64
	src    CellID
}

func (self *rpcResponse) Feedback(data interface{}) {

	RawSend(self.src, data, self.callid)
}

func (self *rpcResponse) GetPacket() *Packet {
	return self.pkt
}

func (self *rpcResponse) GetSession() CellID {
	return self.src
}

type RPCResponse interface {
	Feedback(interface{})
	GetPacket() *Packet
}

// rpc每次调用上下文
type remoteCall struct {
	Done chan bool

	Reply interface{} // 用户返回结构体

	callid int64
	e      error
}

func (self *remoteCall) done() {
	self.Done <- true
}

func (self *remoteCall) fail(e error) {
	self.e = e
	self.Done <- false
}

var (
	callMap      = make(map[int64]*remoteCall)
	callIDAcc    int64
	callMapGuard sync.Mutex
)

func addCall(c *remoteCall) int64 {

	callMapGuard.Lock()

	defer callMapGuard.Unlock()

	callIDAcc++

	callid := callIDAcc

	callMap[callid] = c

	return callid
}

func removeCall(id int64) {
	callMapGuard.Lock()

	delete(callMap, id)

	callMapGuard.Unlock()
}

func getCall(id int64) *remoteCall {
	callMapGuard.Lock()

	defer callMapGuard.Unlock()

	if v, ok := callMap[id]; ok {
		return v
	}

	return nil
}

var (
	errRequestTimeout error         = errors.New("RPC reqest time out")
	timeOut           time.Duration = time.Second * 5
)

type localData struct {
	Session CellID
	Packet  *Packet
}

func (self localData) GetSession() CellID {
	return self.Session
}

func (self localData) GetPacket() *Packet {
	return self.Packet
}

func InjectPost(target CellID, data interface{}, callid int64, src CellID) {

	pkt := data.(*Packet)

	if callid != 0 {

		LocalPost(target, &rpcResponse{
			pkt: pkt,

			callid: callid,
			src:    src,
		})

	} else {
		LocalPost(target, &localData{Packet: pkt, Session: src})
	}

}

func Call(target CellID, data interface{}, src CellID) (interface{}, error) {

	c := &remoteCall{Done: make(chan bool)}

	c.callid = addCall(c)

	// 投递内容就在本地, 马上post
	if IsLocal(target) {

		InjectPost(target, data, c.callid, src)

	} else {

		// 真正的rpc
		ExpressPost(target, data, c.callid, src)

	}

	select {
	// 等待异步响应
	case <-c.Done:
		return c.Reply, nil
	case <-time.After(timeOut):
		removeCall(c.callid)
		return nil, errRequestTimeout
	}

	return nil, nil
}
