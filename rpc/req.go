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

	recv, send := p.GetHandler()

	installSendHandler(p, send)
	installRecvHandler(p, recv, args, userCallback)

	ses.Send(args)
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
