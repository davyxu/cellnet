package rpc

import (
	"errors"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

var (
	errInvalidPeerSession   error = errors.New("rpc: invalid peer type, require connector")
	errConnectorSesNotReady error = errors.New("rpc: connector session not ready")
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(ud interface{}, reqMsg interface{}, userCallback interface{}) error {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return err
	}

	recv, send := p.GetHandler()

	installSendHandler(p, send)

	if err := installAsyncRecvHandler(p, recv, reqMsg, userCallback); err != nil {
		return err
	}

	ses.Send(reqMsg)

	return nil
}

// 发出请求, 接收到服务器返回后才返回, ud: peer/session,   reqMsg:请求用的消息, ackMsgName: 返回消息类型名, 返回消息为返回值
func CallSync(ud interface{}, reqMsg interface{}, ackMsgName string) (interface{}, error) {

	ses, p, err := getPeerSession(ud)

	if err != nil {
		return nil, err
	}

	recv, send := p.GetHandler()

	installSendHandler(p, send)

	ret := make(chan interface{})

	if err := installSyncRecvHandler(p, recv, reqMsg, ackMsgName, ret); err != nil {
		return nil, err
	}

	ses.Send(reqMsg)

	return <-ret, nil
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
