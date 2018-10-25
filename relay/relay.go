package relay

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

var (
	ErrInvalidPeerSession = errors.New("Require valid cellnet.Session or cellnet.TCPConnector")
)

// sesDetector: 提供要发送到的目标session， 传输msg消息，透传passThroughData
func Relay(sesDetector interface{}, payloadList ...interface{}) error {

	ses, err := getSession(sesDetector)
	if err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	var ack RelayACK

	for _, payload := range payloadList {
		switch value := payload.(type) {
		case int64:
			ack.Int64 = value
			ack.Type = RelayPassThroughType_Int64
		case []int64:
			ack.Int64Slice = value
			ack.Type = RelayPassThroughType_Int64Slice
		case []byte:
			ack.Bytes = value
		default:
			var meta *cellnet.MessageMeta
			ack.Msg, meta, err = codec.EncodeMessage(payload, nil)

			if err != nil {
				return err
			}

			ack.MsgID = uint32(meta.ID)
		}
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
