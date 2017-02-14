package cellnet

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 消息到封包
func BuildPacket(msgObj interface{}) ([]byte, *MessageMeta) {

	msg := msgObj.(proto.Message)

	rawdata, err := proto.Marshal(msg)

	if err != nil {
		log.Errorln(err)
	}

	meta := MessageMetaByName(MessageFullName(reflect.TypeOf(msg)))

	return rawdata, meta
}

// 封包到消息
func ParsePacket(data []byte, msgType reflect.Type) (interface{}, error) {
	// msgType 为ptr类型, new时需要非ptr型

	rawMsg := reflect.New(msgType.Elem()).Interface()

	err := proto.Unmarshal(data, rawMsg.(proto.Message))

	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}
