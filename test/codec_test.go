package tests

import (
	"fmt"
	cellcodec "github.com/davyxu/cellnet/codec"
	cellmeta "github.com/davyxu/cellnet/meta"
	"github.com/davyxu/x/container"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type CodecPanicMsg struct {
}

func (self *CodecPanicMsg) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellmeta.Register(&cellmeta.Meta{
		Codec: new(problemJsonCodec),
		Type:  reflect.TypeOf((*CodecPanicMsg)(nil)).Elem(),
		Id:    1111,
	})
}

type problemJsonCodec struct {
}

func (self *problemJsonCodec) Name() string {
	return "problemjson"
}

// 将结构体编码为JSON的字节数组
func (self *problemJsonCodec) Encode(msgObj any, ps *xcontainer.Mapper) (data any, err error) {
	panic("encode")
	return nil, nil
}

func (self *problemJsonCodec) Decode(data any, msgObj any) error {
	panic("decode")
	return nil
}

func TestPanic(t *testing.T) {

	assert.NotPanics(t, func() {

		_, _, err := cellcodec.Encode(&CodecPanicMsg{}, nil)

		assert.Equal(t, err.Error(), "encode panic: encode")
	})

	assert.NotPanics(t, func() {

		_, _, err := cellcodec.Decode(1111, nil)

		assert.Equal(t, err.Error(), "encode panic: decode")
	})

}
