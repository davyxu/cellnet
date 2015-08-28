package cellnet

import (
	"errors"
	"sync"
	"time"
)

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

func Call(target CellID, data interface{}) (interface{}, error) {

	c := &remoteCall{Done: make(chan bool)}

	c.callid = addCall(c)

	// TODO 可能在这里卡死影响后面的等待
	RawSend(target, data, c.callid)

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
