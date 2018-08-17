package relay

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type RelayACK struct {
	PayloadMsgID uint16
	Payload      []byte // 用户原消息

	PassThroughKind  uint8
	PassThroughMsgID uint16
	PassThrough      []byte // 透传数据
}

var (
	ErrRelayPacketCrack = errors.New("invalid relay packet format")
)

func (self *RelayACK) Decode() (payload, passThrough interface{}, err error) {

	payload, _, err = codec.DecodeMessage(int(self.PayloadMsgID), self.Payload)
	if err != nil {
		return
	}

	switch self.PassThroughKind {
	case 1:
		passThrough = int64(binary.LittleEndian.Uint64(self.PassThrough))
	case 2:
		passThrough, _, err = codec.DecodeMessage(int(self.PassThroughMsgID), self.PassThrough)

	default:
		err = ErrRelayPacketCrack
	}

	return
}

func (self *RelayACK) Encode(payload, passThrough interface{}) (err error) {
	var payloadMeta *cellnet.MessageMeta

	self.Payload, payloadMeta, err = codec.EncodeMessage(payload, nil)

	if err != nil {
		return
	}

	self.PayloadMsgID = uint16(payloadMeta.ID)

	var passThroughMeta *cellnet.MessageMeta

	switch ptValue := passThrough.(type) {
	case int64:

		self.PassThrough = make([]byte, 8)

		binary.LittleEndian.PutUint64(self.PassThrough, uint64(ptValue))
		self.PassThroughKind = 1

	default:
		self.PassThrough, passThroughMeta, err = codec.EncodeMessage(passThrough, nil)

		if err != nil {
			return
		}

		self.PassThroughMsgID = uint16(passThroughMeta.ID)
		self.PassThroughKind = 2
	}

	return
}

func (self *RelayACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*RelayACK)(nil)).Elem(),
		ID:    int(util.StringHash("relay.RelayACK")),
	})

}
