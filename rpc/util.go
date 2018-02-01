package rpc

import (
	"errors"
	"github.com/davyxu/cellnet"
)

var (
	ErrInvalidPeerSession = errors.New("rpc: Invalid peer type, require cellnet.RPCSessionGetter or cellnet.Session")
)

type RPCSessionGetter interface {
	RPCSession() cellnet.Session
}

// 从peer获取rpc使用的session
func getPeerSession(ud interface{}) (cellnet.Session, error) {

	if ud == nil {
		return nil, ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		return i.RPCSession(), nil
	case cellnet.Session:
		return i, nil
	default:
		return nil, ErrInvalidPeerSession
	}
}
