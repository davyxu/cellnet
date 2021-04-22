package tests

import (
	"fmt"
	cellcodec "github.com/davyxu/cellnet/codec"
	cellmeta "github.com/davyxu/cellnet/meta"
	xframe "github.com/davyxu/x/frame"
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
		ID:    1111,
	})
}

type problemJsonCodec struct {
}

func (self *problemJsonCodec) Name() string {
	return "problemjson"
}

// 将结构体编码为JSON的字节数组
func (self *problemJsonCodec) Encode(msgObj interface{}, ps *xframe.PropertySet) (data interface{}, err error) {
	panic("encode")
	return nil, nil
}

func (self *problemJsonCodec) Decode(data interface{}, msgObj interface{}) error {
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
