package rpc

import (
	"errors"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
)

var (
	ErrInvalidPeerSession    error = errors.New("rpc: Invalid peer type, require connector")
	ErrConnectorSesNotReady  error = errors.New("rpc: Connector session not ready")
	ErrReplayMessageNotFound       = errors.New("rpc: Reply message name not found")

	metaWrapper = cellnet.MessageMetaByName("gamedef.RemoteCallACK")
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

			return nil, nil, ErrInvalidPeerSession
		}
	case cellnet.Session:
		ses = ud.(cellnet.Session)
	}

	if ses == nil {
		return nil, nil, ErrConnectorSesNotReady
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
