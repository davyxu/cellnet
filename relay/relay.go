package relay

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

var (
	ErrInvalidPeerSession = errors.New("relay: Require cellnet.Session")
)

// sesDetector: 提供要发送到的目标session， 发送msg消息，并携带ContextID
func Relay(sesDetector, msg interface{}, contextIDList ...int64) error {

	ses, err := getSession(sesDetector)
	if err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	data, meta, err := codec.EncodeMessage(msg)

	if err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	ses.Send(&RelayACK{
		MsgID:     uint16(meta.ID),
		Data:      data,
		ContextID: contextIDList,
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
