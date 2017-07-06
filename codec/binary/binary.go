package binary

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/goobjfmt"
)

type binaryCodec struct {
}

func (self *binaryCodec) Name() string {
	return "binary"
}

func (self *binaryCodec) Encode(msgObj interface{}) ([]byte, error) {

	return goobjfmt.BinaryWrite(msgObj)

}

func (self *binaryCodec) Decode(data []byte, msgObj interface{}) error {

	return goobjfmt.BinaryRead(data, msgObj)
}

func init() {

	cellnet.RegisterCodec("binary", new(binaryCodec))
}
