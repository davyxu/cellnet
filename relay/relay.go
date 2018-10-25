package relay

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

var (
	ErrInvalidPeerSession = errors.New("Require valid cellnet.Session or cellnet.TCPConnector")
)

// payload: msg/bytes   passthrough: int64, []int64, string
func Relay(sesDetector interface{}, dataList ...interface{}) error {

	ses, err := getSession(sesDetector)
	if err != nil {
		log.Errorln("relay.Relay:", err)
		return err
	}

	var ack RelayACK

	for _, payload := range dataList {
		switch value := payload.(type) {
		case int64:
			ack.Int64 = value
		case []int64:
			ack.Int64Slice = value

		case string:
			ack.Str = value
		case []byte: // 作为payload
			ack.Bytes = value
		default:
			if ack.MsgID == 0 {
				var meta *cellnet.MessageMeta
				ack.Msg, meta, err = codec.EncodeMessage(payload, nil)

				if err != nil {
					return err
				}

				ack.MsgID = uint32(meta.ID)
			} else {
				panic("Multi message relay not support")
			}

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
