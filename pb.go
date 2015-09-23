package cellnet

import (
	"github.com/davyxu/cellnet/util"
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

func Name2ID(name string) int {

	return int(util.StringHashNoCase(name))
}

func Type2ID(msg proto.Message) int {

	name := reflect.TypeOf(msg).Elem().String()

	return int(util.StringHashNoCase(name))
}

func ReflectProtoName(msg proto.Message) string {

	return reflect.TypeOf(msg).Elem().String()
}

// TODO 抽象Filter

// 消息到封包
func BuildPacket(data interface{}) *Packet {

	msg := data.(proto.Message)

	rawdata, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal(err)
	}

	msgID := uint32(Name2ID(ReflectProtoName(msg)))

	return &Packet{
		MsgID: msgID,
		Data:  rawdata,
	}
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
