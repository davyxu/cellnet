package cellnet

import (
	"reflect"
)

type MessageMeta struct {
	Type reflect.Type
	Name string
	ID   int
}

func NewMessageMeta(msg interface{}) *MessageMeta {

	msgType := reflect.TypeOf(msg)

	msgName := msgType.String()

	return &MessageMeta{
		Type: msgType,
		Name: msgName,
		ID:   Name2ID(msgName),
	}
}
