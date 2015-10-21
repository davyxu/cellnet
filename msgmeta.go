package cellnet

import (
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   int
}

func NewMessageMeta(msg interface{}) *MessageMeta {

	msgType := reflect.TypeOf(msg)

	if msgType.Kind() == reflect.Ptr {
		msgType = msgType.Elem()
	}

	msgName := msgType.String()

	return &MessageMeta{
		Type: msgType,
		Name: msgName,
		ID:   int(util.StringHashNoCase(msgName)),
	}
}
