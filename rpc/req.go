package rpc

import (
	"errors"
	"reflect"
	"sync"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/gamedef"
	"github.com/davyxu/cellnet/socket"
)

var (
	reqByID  = make(map[int64]*request)
	reqGuard sync.RWMutex
	idacc    int64
)

var (
	errInvalidPeerSession   error = errors.New("rpc: invalid peer type, require connector")
	errConnectorSesNotReady error = errors.New("rpc: connector session not ready")
)

// 添加一个rpc的调用信息
func addCall(req *request) {

	reqGuard.Lock()

	idacc++

	req.id = idacc

	// TODO 底层支持timer, 抛出一个超时检测, 清理map
	reqByID[req.id] = req

	reqGuard.Unlock()

}

// 获取一个rpc调用信息
func getCall(id int64) *request {

	reqGuard.RLock()
	defer reqGuard.RUnlock()

	if v, ok := reqByID[id]; ok {

		return v
	}

	return nil
}

func removeCall(id int64) {
	reqGuard.Lock()
	delete(reqByID, id)
	reqGuard.Unlock()
}

// 从peer获取rpc使用的session
func getPeerSession(p interface{}) (cellnet.Session, cellnet.EventQueue, error) {

	var ses cellnet.Session

	switch p.(type) {
	case cellnet.Peer:
		if connPeer, ok := p.(interface {
			DefaultSession() cellnet.Session
		}); ok {

			ses = connPeer.DefaultSession()

		} else {

			return nil, nil, errInvalidPeerSession
		}
	case cellnet.Session:
		ses = p.(cellnet.Session)
	}

	if ses == nil {
		return nil, nil, errConnectorSesNotReady
	}

	return ses, ses.FromPeer().(cellnet.EventQueue), nil
}

// 传入peer或者session
func Call(p interface{}, args interface{}, callback interface{}) {

	ses, evq, err := getPeerSession(p)

	if err != nil {
		log.Errorln(err)
		return
	}

	_, msg := newRequest(evq, args, callback)

	ses.Send(msg)

	// TODO rpc日志
}

// 传入peer或者session
func CallSync(p interface{}, args interface{}, callback interface{}) {
	ses, evq, err := getPeerSession(p)

	if err != nil {
		log.Errorln(err)
		return
	}

	req, msg := newRequest(evq, args, callback)
	req.recvied = make(chan bool)

	ses.Send(msg)

	<-req.recvied

}

type request struct {
	id        int64
	callback  reflect.Value
	replyType reflect.Type
	recvied   chan bool
}

func (self *request) done(msg *gamedef.RemoteCallACK) {

	rawType, err := cellnet.ParsePacket(&cellnet.Packet{
		MsgID: msg.MsgID,
		Data:  msg.Data,
	}, self.replyType)

	defer removeCall(self.id)

	if err != nil {
		log.Errorln(err)
		return
	}

	// 这里的反射, 会影响非常少的效率, 但因为外部写法简单, 就算了
	self.callback.Call([]reflect.Value{reflect.ValueOf(rawType)})

	if self.recvied != nil {
		self.recvied <- true
	}

}

var needRegisterClient bool = true
var needRegisterClientGuard sync.Mutex

func newRequest(evq cellnet.EventQueue, args interface{}, callback interface{}) (*request, interface{}) {

	needRegisterClientGuard.Lock()
	if needRegisterClient {
		// 请求端
		socket.RegisterSessionMessage(evq, "gamedef.RemoteCallACK", func(content interface{}, ses cellnet.Session) {
			msg := content.(*gamedef.RemoteCallACK)

			c := getCall(msg.CallID)

			if c == nil {
				return
			}

			c.done(msg)
		})

		needRegisterClient = false
	}
	needRegisterClientGuard.Unlock()

	req := &request{}

	funcType := reflect.TypeOf(callback)
	req.replyType = funcType.In(0)
	req.callback = reflect.ValueOf(callback)

	pkt, _ := cellnet.BuildPacket(args)

	addCall(req)

	return req, &gamedef.RemoteCallREQ{
		MsgID:  pkt.MsgID,
		Data:   pkt.Data,
		CallID: req.id,
	}
}
