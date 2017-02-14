package rpc

import (
	"errors"
	"reflect"
	"sync"

	"github.com/davyxu/cellnet"
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
func getPeerSession(ud interface{}) (cellnet.Session, cellnet.Peer, error) {

	var ses cellnet.Session

	switch ud.(type) {
	case cellnet.Peer:
		if connPeer, ok := ud.(interface {
			DefaultSession() cellnet.Session
		}); ok {

			ses = connPeer.DefaultSession()

		} else {

			return nil, nil, errInvalidPeerSession
		}
	case cellnet.Session:
		ses = ud.(cellnet.Session)
	}

	if ses == nil {
		return nil, nil, errConnectorSesNotReady
	}

	return ses, ses.FromPeer(), nil
}

// 传入peer或者session
func Call(ud interface{}, args interface{}, userCallback func(*cellnet.SessionEvent)) {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		log.Errorln(err)
		return
	}

	newRequest(p, args, userCallback)

	ses.Send(args)
}

// 传入peer或者session
//func CallSync(p interface{}, args interface{}, callback interface{}) {
//	ses, evq, err := getPeerSession(p)

//	if err != nil {
//		log.Errorln(err)
//		return
//	}

//	req, msg := newRequest(evq, args, callback)
//	req.recvied = make(chan bool)

//	ses.Send(msg)

//	<-req.recvied

//}

type request struct {
	id        int64
	callback  reflect.Value
	replyType reflect.Type
	recvied   chan bool
}

// socket.EncodePacketHandler -> socket.MsgLogHandler -> rpc.BoxHandler -> socket.WritePacketHandler

func installSendHandler(p cellnet.Peer, send cellnet.EventHandler) {
	// 发送的Handler
	if cellnet.HandlerName(send) == "EncodePacketHandler" {

		var start cellnet.EventHandler

		if cellnet.HandlerName(send.Next()) == "MsgLogHandler" {
			start = send.Next()
		} else {
			start = send
		}

		// 已经装过了
		if start.MatchTag("rpc") {
			return
		}

		first := NewBoxHandler()
		first.SetTag("rpc")

		cellnet.LinkHandler(start, first, socket.NewWritePacketHandler())

	} else {
		panic("unknown send handler struct")
	}
}

//  socket.DispatcherHandler -> rpc.UnboxHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func installRecvHandler(p cellnet.Peer, recv cellnet.EventHandler, args interface{}, userCallback func(*cellnet.SessionEvent)) {

	// 接收
	msgName := cellnet.MessageFullName(reflect.TypeOf(args))
	meta := cellnet.MessageMetaByName(msgName)

	if meta == nil {
		panic("can not found rpc message:" + msgName)
	}

	// RPC消息只能被注册1个
	if p.CountByID(int(meta.ID)) == 0 {

		p.AddHandler(int(meta.ID), buildRecvHandler(meta, userCallback, nil))

	}

}

func newRequest(p cellnet.Peer, args interface{}, userCallback func(*cellnet.SessionEvent)) {

	recv, send := p.GetHandler()

	installSendHandler(p, send)
	installRecvHandler(p, recv, args, userCallback)

}
