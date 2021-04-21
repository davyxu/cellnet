package tests

import (
	"fmt"
	cellcodec "github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/json"
	cellmeta "github.com/davyxu/cellnet/meta"
	xbytes "github.com/davyxu/x/bytes"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellmeta.Register(&cellmeta.Meta{
		Codec: cellcodec.MustGetByName("json"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(xbytes.StringHash("TestEchoACK")),
	})
}
