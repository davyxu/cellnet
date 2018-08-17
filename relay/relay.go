package relay

import (
	"errors"
	"github.com/davyxu/cellnet"
)

var (
	ErrInvalidPeerSession = errors.New("Require valid cellnet.Session or cellnet.TCPConnector")
)

// sesDetector: 提供要发送到的目标session， 传输msg消息，透传passThroughData
func Relay(sesDetector, payload, passThrough interface{}) error {

	ses, err := getSession(sesDetector)
	if err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	var ack RelayACK
	if err = ack.Encode(payload, passThrough); err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	ses.Send(&ack)

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
