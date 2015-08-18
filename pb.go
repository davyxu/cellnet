package cellnet

import (
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

func Name2ID(name string) int {

	return int(StringHashNoCase(name))
}

func Type2ID(msg proto.Message) int {

	name := reflect.TypeOf(msg).Elem().String()

	return int(StringHashNoCase(name))
}

func ReflectProtoName(msg proto.Message) string {

	return reflect.TypeOf(msg).Elem().String()
}

// 消息到封包
func BuildPacket(msg proto.Message) *Packet {

	data, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal(err)
	}

	msgID := uint32(Name2ID(ReflectProtoName(msg)))

	return &Packet{
		MsgID: msgID,
		Data:  data,
	}
}

func ParsePacket(pkt *Packet, msgType reflect.Type) (interface{}, error) {

	rawMsg := reflect.New(msgType).Interface()

	err := proto.Unmarshal(pkt.Data, rawMsg.(proto.Message))

	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}
