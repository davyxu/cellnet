package proto

import (
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/binary"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta("binary", // 消息的编码格式
		"test.TestEchoACK",                         // 消息名
		reflect.TypeOf((*TestEchoACK)(nil)).Elem(), // 消息的反射类型
		1, // 消息ID
	)
}
