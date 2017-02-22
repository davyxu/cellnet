package rpc

import (
	"errors"

	"github.com/davyxu/cellnet"
)

var (
	ErrInvalidPeerSession    error = errors.New("rpc: Invalid peer type, require connector")
	ErrConnectorSesNotReady  error = errors.New("rpc: Connector session not ready")
	ErrReplayMessageNotFound       = errors.New("rpc: Reply message name not found")

	metaWrapper = cellnet.MessageMetaByName("coredef.RemoteCallACK")
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
