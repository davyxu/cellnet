package proto

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/binary"
	"reflect"
)

type ChatREQ struct {
	Content string
}

type ChatACK struct {
	Content string
	Id      int64
}

func (self *ChatREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *ChatACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*ChatREQ)(nil)).Elem(),
		ID:    501,
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*ChatACK)(nil)).Elem(),
		ID:    502,
	})
}
