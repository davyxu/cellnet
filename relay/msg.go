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
	ErrRelayPacketCrack       = errors.New("invalid relay packet format")
	ErrUnknownPassThroughKind = errors.New("Unknown PassThrough kind")
)

func (self *RelayACK) Decode() (payload, passThrough interface{}, err error) {

	payload, _, err = codec.DecodeMessage(int(self.PayloadMsgID), self.Payload)
	if err != nil {
		return
	}

	switch self.PassThroughKind {
	case 0: // msg
		passThrough, _, err = codec.DecodeMessage(int(self.PassThroughMsgID), self.PassThrough)
	case 1: // int64
		passThrough = int64(binary.LittleEndian.Uint64(self.PassThrough))
	case 2: // []int64
		if len(self.PassThrough) >= 2 {
			ptr := self.PassThrough
			dataLen := binary.LittleEndian.Uint16(ptr)
			ptr = ptr[2:]

			list := make([]int64, dataLen)

			for i := uint16(0); i < dataLen; i++ {
				list[i] = int64(binary.LittleEndian.Uint64(ptr))

				ptr = ptr[8:]
			}

			passThrough = list
		} else {
			passThrough = make([]int64, 0)
		}

	case 3:
		passThrough = nil

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
		self.PassThroughKind = 1

		self.PassThrough = make([]byte, 8)

		binary.LittleEndian.PutUint64(self.PassThrough, uint64(ptValue))

	case []int64:
		self.PassThroughKind = 2

		self.PassThrough = make([]byte, 8*len(ptValue)+2)
		ptr := self.PassThrough

		binary.LittleEndian.PutUint16(ptr, uint16(len(ptValue)))

		ptr = ptr[2:]

		for _, v := range ptValue {
			binary.LittleEndian.PutUint64(ptr, uint64(v))
			ptr = ptr[8:]
		}

	case nil:
		self.PassThroughKind = 3

	default:

		ptType := reflect.TypeOf(passThrough)
		if ptType.Kind() == reflect.Ptr && ptType.Elem().Kind() == reflect.Struct {
			self.PassThrough, passThroughMeta, err = codec.EncodeMessage(passThrough, nil)

			if err != nil {
				return
			}

			self.PassThroughMsgID = uint16(passThroughMeta.ID)
			self.PassThroughKind = 0
		} else {
			err = ErrUnknownPassThroughKind
		}
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
