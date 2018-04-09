package rpc

import (
	"errors"
	"github.com/davyxu/cellnet"
)

var (
	ErrInvalidPeerSession = errors.New("rpc: Invalid peer type, require cellnet.RPCSessionGetter or cellnet.Session")
	ErrEmptySession       = errors.New("rpc: Empty session")
)

type RPCSessionGetter interface {
	RPCSession() cellnet.Session
}

// 从peer获取rpc使用的session
func getPeerSession(ud interface{}) (ses cellnet.Session, err error) {

	if ud == nil {
		return nil, ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		ses = i.RPCSession()
	case cellnet.Session:
		ses = i
	case cellnet.TCPConnector:
		ses = i.Session()
	default:
		err = ErrInvalidPeerSession
		return
	}

	if ses == nil {
		return nil, ErrEmptySession
	}

	return
}
