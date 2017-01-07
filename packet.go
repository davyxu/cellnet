package cellnet

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 普通封包
type Packet struct {
	MsgID uint32 // 消息ID
	Data  []byte
}

func (self Packet) ContextID() uint32 {
	return self.MsgID
}

// 消息到封包
func BuildPacket(data interface{}) (*Packet, *MessageMeta) {

	msg := data.(proto.Message)

	rawdata, err := proto.Marshal(msg)

	if err != nil {
		log.Errorln(err)
	}

	meta := MessageMetaByType(reflect.TypeOf(msg))

	return &Packet{
		MsgID: uint32(meta.ID),
		Data:  rawdata,
	}, meta
}

// 封包到消息
func ParsePacket(pkt *Packet, msgType reflect.Type) (interface{}, error) {
	// msgType 为ptr类型, new时需要非ptr型

	rawMsg := reflect.New(msgType.Elem()).Interface()

	err := proto.Unmarshal(pkt.Data, rawMsg.(proto.Message))

	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}
