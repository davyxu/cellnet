package cellnet

import (
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

// TODO 抽象Filter

// 消息到封包
func BuildPacket(data interface{}) (*Packet, *MessageMeta) {

	msg := data.(proto.Message)

	rawdata, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal(err)
	}

	meta := NewMessageMeta(msg)

	return &Packet{
		MsgID: uint32(meta.ID),
		Data:  rawdata,
	}, meta
}

// 封包到消息
func ParsePacket(pkt *Packet, msgType reflect.Type) (interface{}, error) {

	rawMsg := reflect.New(msgType).Interface()

	err := proto.Unmarshal(pkt.Data, rawMsg.(proto.Message))

	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}
