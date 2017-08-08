package rpc

import (
	"errors"

	"github.com/davyxu/cellnet"
)

var (
	ErrInvalidPeerSession    error = errors.New("rpc: Invalid peer type, require cellnet.RPCSessionGetter or cellnet.Session")
	ErrReplayMessageNotFound       = errors.New("rpc: Reply message name not found")

	MetaCall = cellnet.MessageMetaByName("coredef.RemoteCallACK")
)

type RPCSessionGetter interface {
	RPCSession() cellnet.Session
}

// 从peer获取rpc使用的session
func getPeerSession(ud interface{}) (cellnet.Session, cellnet.Peer, error) {

	if ud == nil {
		return nil, nil, ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		return i.RPCSession(), i.RPCSession().FromPeer(), nil
	case cellnet.Session:
		return i, i.FromPeer(), nil
	default:
		return nil, nil, ErrInvalidPeerSession
	}

}
