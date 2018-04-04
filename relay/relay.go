package relay

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

var (
	ErrInvalidPeerSession = errors.New("relay: Require cellnet.Session")
)

func Relay(sesDetector, msg interface{}, sesid ...int64) error {

	ses, err := getSession(sesDetector)
	if err != nil {
		return err
	}

	data, meta, err := codec.EncodeMessage(msg)

	if err != nil {
		return err
	}

	ses.Send(&RelayACK{
		MsgID:     uint16(meta.ID),
		Data:      data,
		SessionID: sesid,
	})

	return nil
}

func getSession(sesDetector interface{}) (cellnet.Session, error) {
	switch unknown := sesDetector.(type) {
	case cellnet.Session:
		return unknown, nil
	case cellnet.TCPConnector:
		return unknown.Session(), nil
	default:
		return nil, ErrInvalidPeerSession
	}
}
