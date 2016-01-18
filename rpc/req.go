package rpc

import (
	"errors"
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
)

var (
	callMap = make(map[int64]*request)
	idacc   int64
)

type connSes interface {
	DefaultSession() cellnet.Session
}

var (
	errInvalidPeerSession   error = errors.New("rpc: invalid peer type, require connector")
	errConnectorSesNotReady error = errors.New("rpc: connector session not ready")
)

// 添加一个rpc的调用信息
func addCall() *request {

	idacc++

	req := &request{
		id: idacc,
	}

	// TODO 底层支持timer, 抛出一个超时检测, 清理map
	callMap[req.id] = req

	return req
}

// 获取一个rpc调用信息
func getCall(id int64) *request {

	if v, ok := callMap[id]; ok {

		return v
	}

	return nil
}

func removeCall(id int64) {
	delete(callMap, id)
}

// 从peer获取rpc使用的session
func getPeerSession(p cellnet.Peer) (cellnet.Session, error) {

	var ses cellnet.Session

	if connPeer, ok := p.(connSes); ok {

		ses = connPeer.DefaultSession()

	} else {

		return nil, errInvalidPeerSession
	}

	if ses == nil {
		return nil, errConnectorSesNotReady
	}

	return ses, nil
}

func Call(p cellnet.Peer, args interface{}, callback interface{}) {

	req := addCall()

	funcType := reflect.TypeOf(callback)
	req.replyType = funcType.In(0).Elem()
	req.callback = reflect.ValueOf(callback)

	ses, err := getPeerSession(p)

	if err != nil {
		log.Errorln(err)
		removeCall(req.id)
		return
	}

	pkt, _ := cellnet.BuildPacket(args)

	ses.Send(&coredef.RemoteCallREQ{
		MsgID:  pkt.MsgID,
		Data:   pkt.Data,
		CallID: req.id,
	})

	// TODO rpc日志

}

type request struct {
	id        int64
	callback  reflect.Value
	replyType reflect.Type
}

func (self *request) done(msg *coredef.RemoteCallACK) {

	rawType, err := cellnet.ParsePacket(&cellnet.Packet{
		MsgID: msg.MsgID,
		Data:  msg.Data,
	}, self.replyType)

	defer removeCall(self.id)

	if err != nil {
		log.Errorln(err)
		return
	}

	self.callback.Call([]reflect.Value{reflect.ValueOf(rawType)})
}

func InstallClient(p cellnet.Peer) {

	// 请求端
	socket.RegisterSessionMessage(p, coredef.RemoteCallACK{}, func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.RemoteCallACK)

		c := getCall(msg.CallID)

		if c == nil {
			return
		}

		c.done(msg)
	})

}
